const { createApp, nextTick } = Vue;

// Fetch configuration from the Go server's /config endpoint.
fetch('/config')
  .then(r => r.json())
  .then(CONFIG => createApp({
  data() {
    return {
      config: CONFIG,
      model: null,            // fetched from /models on startup
      modelError: null,       // error message if /models fetch fails
      messages: [],           // { role, content, rawContent, reasoning, images[], isStreaming, cancelled, friendlyError, errorDetail, reasoningOpen }
      userInput: '',
      isLoading: false,
      showSettings: false,
      attachedImage: null,    // base64 data URL of the image to attach
      lightboxImage: null,    // image URL to show full-size
      abortController: null,  // AbortController for the active request
      errorModal: null,       // raw error text shown in the error detail modal
    };
  },

  mounted() {
    document.title = this.uiTitle;
    window.addEventListener('keydown', this.handleGlobalKeydown);
    this.fetchModel();
  },

  beforeUnmount() {
    window.removeEventListener('keydown', this.handleGlobalKeydown);
    if (this.abortController) this.abortController.abort();
  },

  computed: {
    uiTitle() {
      const name = this.config.instanceName ?? '';
      return name.charAt(0).toUpperCase() + name.slice(1) + ' Inference Snap';
    },
    supportsVision() {
      return this.config.capabilities?.includes('vision') ?? false;
    },
  },

  methods: {
    // ----------------------------------------------------------------
    // Model discovery
    // ----------------------------------------------------------------
    async fetchModel() {
      try {
        const response = await fetch(`${this.config.openAIBaseURL}/models`);
        if (!response.ok) throw new Error(`${response.status} ${response.statusText}`);
        const data = await response.json();
        const first = data.data?.[0];
        if (!first) throw new Error('No models returned by the API.');
        this.model = first.id;
      } catch (err) {
        this.modelError = err.message;
      }
    },

    // ----------------------------------------------------------------
    // Keyboard / lightbox / error modal
    // ----------------------------------------------------------------
    handleGlobalKeydown(e) {
      if (e.key === 'Escape') {
        if (this.lightboxImage) this.lightboxImage = null;
        if (this.errorModal) this.errorModal = null;
      }
    },

    openLightbox(src) {
      this.lightboxImage = src;
    },

    closeLightbox() {
      this.lightboxImage = null;
    },

    openErrorModal(detail) {
      this.errorModal = detail;
    },

    closeErrorModal() {
      this.errorModal = null;
    },

    // ----------------------------------------------------------------
    // Text formatting (HTML-safe)
    // ----------------------------------------------------------------
    formatText(text) {
      if (!text) return '';
      return String(text)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;')
        .replace(/\n/g, '<br>');
    },

    // ----------------------------------------------------------------
    // Input handling
    // ----------------------------------------------------------------
    handleKeydown(event) {
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault();
        this.sendMessage();
      }
      // Shift+Enter falls through to the textarea's default behavior (inserts newline).
    },

    // ----------------------------------------------------------------
    // Image handling
    // ----------------------------------------------------------------
    resizeImage(file) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onerror = () => reject(new Error('Failed to read image file.'));
        reader.onload = (e) => {
          const img = new Image();
          img.onerror = () => reject(new Error('Failed to decode image.'));
          img.onload = () => {
            const MAX_PIXELS = 500_000; // 0.5 megapixels
            const pixels = img.width * img.height;
            let width = img.width;
            let height = img.height;
            if (pixels > MAX_PIXELS) {
              const ratio = Math.sqrt(MAX_PIXELS / pixels);
              width = Math.floor(img.width * ratio);
              height = Math.floor(img.height * ratio);
            }
            const canvas = document.createElement('canvas');
            canvas.width = width;
            canvas.height = height;
            canvas.getContext('2d').drawImage(img, 0, 0, width, height);
            resolve(canvas.toDataURL('image/jpeg', 0.85));
          };
          img.src = e.target.result;
        };
        reader.readAsDataURL(file);
      });
    },

    async handleImageUpload(event) {
      const file = event.target.files[0];
      if (!file) return;
      if (!file.type.startsWith('image/')) {
        alert('Only image files are supported.');
        event.target.value = '';
        return;
      }
      try {
        this.attachedImage = await this.resizeImage(file);
      } catch (err) {
        alert(`Could not load image: ${err.message}`);
      }
      event.target.value = '';
    },

    removeImage() {
      this.attachedImage = null;
    },

    // ----------------------------------------------------------------
    // Reasoning: parse <think>…</think> tags into a separate field
    // ----------------------------------------------------------------
    syncReasoningFromRaw(msg) {
      const raw = msg.rawContent;
      const thinkMatch = /<think>([\s\S]*?)(<\/think>|$)/.exec(raw);
      if (thinkMatch) {
        msg.reasoning = thinkMatch[1];
        msg.content = raw.replace(/<think>[\s\S]*?(<\/think>|$)/g, '').trim();
      } else {
        msg.content = raw;
      }
    },

    // ----------------------------------------------------------------
    // Stop generation
    // ----------------------------------------------------------------
    stopGeneration() {
      if (this.abortController) {
        this.abortController.abort();
      }
    },

    // ----------------------------------------------------------------
    // Retry after error: remove the failed assistant message and
    // re-send the original user message.
    // ----------------------------------------------------------------
    retryMessage(assistantMsgIndex) {
      const userMsgIndex = assistantMsgIndex - 1;
      const userMsg = this.messages[userMsgIndex];
      if (!userMsg || userMsg.role !== 'user') return;

      const content = userMsg.content;
      this.userInput = typeof content === 'string'
        ? content
        : (content.find(p => p.type === 'text')?.text ?? '');

      if (userMsg.images?.length) {
        this.attachedImage = userMsg.images[0];
      }

      this.messages.splice(userMsgIndex, 2);
      this.sendMessage();
    },

    // ----------------------------------------------------------------
    // Send message
    // ----------------------------------------------------------------
    async sendMessage() {
      if (this.isLoading) return;
      if (!this.model) {
        alert('Model is not available yet. Please wait or check the server logs.');
        return;
      }
      const text = this.userInput.trim();
      if (!text && !this.attachedImage) return;

      // Build content for the API.
      let apiContent;
      if (this.attachedImage && this.supportsVision) {
        apiContent = [
          { type: 'text', text: text || '(image attached)' },
          { type: 'image_url', image_url: { url: this.attachedImage } },
        ];
      } else {
        apiContent = text;
      }

      this.messages.push({
        role: 'user',
        content: apiContent,
        rawContent: '',
        reasoning: '',
        images: this.attachedImage ? [this.attachedImage] : [],
        isStreaming: false,
      });

      this.userInput = '';
      this.attachedImage = null;
      this.isLoading = true;

      await nextTick();
      this.scrollToBottom(true);

      // Placeholder for the assistant reply (filled by the stream).
      this.messages.push({
        role: 'assistant',
        content: '',
        rawContent: '',
        reasoning: '',
        images: [],
        isStreaming: true,
        cancelled: false,
        friendlyError: null,
        errorDetail: null,
        reasoningOpen: true,
      });
      const msgIndex = this.messages.length - 1;

      this.abortController = new AbortController();

      try {
        const apiMessages = this.messages
          .filter(m => !m.isStreaming && !m.cancelled)
          .map(({ role, content }) => ({ role, content }));

        const response = await fetch(`${this.config.openAIBaseURL}/chat/completions`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          signal: this.abortController.signal,
          body: JSON.stringify({
            model: this.model,
            messages: apiMessages,
            stream: true,
          }),
        });

        if (!response.ok) {
          const errorText = await response.text().catch(() => '');
          throw new Error(
            `API error ${response.status}: ${response.statusText}` +
            (errorText ? '\n\n' + errorText : '')
          );
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let buffer = '';

        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          buffer += decoder.decode(value, { stream: true });
          const lines = buffer.split('\n');
          buffer = lines.pop();

          for (const line of lines) {
            if (!line.startsWith('data: ')) continue;
            const data = line.slice(6).trim();
            if (data === '[DONE]') continue;
            if (!this.messages[msgIndex]) break;
            try {
              const chunk = JSON.parse(data);
              const delta = chunk.choices?.[0]?.delta ?? {};
              if (delta.reasoning_content) {
                this.messages[msgIndex].reasoning += delta.reasoning_content;
              }
              if (delta.content) {
                this.messages[msgIndex].rawContent += delta.content;
              }
            } catch (_) {}
          }

          this.syncReasoningFromRaw(this.messages[msgIndex]);
          await nextTick();
          this.scrollToBottom();
        }

        this.syncReasoningFromRaw(this.messages[msgIndex]);

      } catch (err) {
        if (err.name === 'AbortError') {
          // Keep partial response visible but exclude it and the triggering user
          // message from future context.
          if (this.messages[msgIndex]) {
            this.messages[msgIndex].cancelled = true;
          }
          const userMsg = this.messages[msgIndex - 1];
          if (userMsg && userMsg.role === 'user') {
            userMsg.cancelled = true;
          }
        } else {
          // Show a friendly error; store raw details for the modal.
          if (this.messages[msgIndex]) {
            this.messages[msgIndex].friendlyError = 'The request failed. Please try again.';
            this.messages[msgIndex].errorDetail = err.message;
          }
        }
      } finally {
        const msg = this.messages[msgIndex];
        if (msg) msg.isStreaming = false;
        this.isLoading = false;
        this.abortController = null;
        await nextTick();
        this.scrollToBottom();
      }
    },

    // ----------------------------------------------------------------
    // Scroll
    // ----------------------------------------------------------------
    scrollToBottom(force = false) {
      const container = this.$refs.messagesContainer;
      if (!container) return;
      const distanceFromBottom = container.scrollHeight - container.scrollTop - container.clientHeight;
      if (force || distanceFromBottom < 100) {
        container.scrollTop = container.scrollHeight;
      }
    },

    toggleReasoning(msg) {
      msg.reasoningOpen = !msg.reasoningOpen;
    },
  },
}).mount('#app'));
