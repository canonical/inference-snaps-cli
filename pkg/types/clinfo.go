package types

type Clinfo struct {
	Platforms []struct {
		ClPlatformName                  string `json:"CL_PLATFORM_NAME"`
		ClPlatformVendor                string `json:"CL_PLATFORM_VENDOR"`
		ClPlatformVersion               string `json:"CL_PLATFORM_VERSION"`
		ClPlatformProfile               string `json:"CL_PLATFORM_PROFILE"`
		ClPlatformExtensions            string `json:"CL_PLATFORM_EXTENSIONS"`
		ClPlatformExtensionsWithVersion struct {
			ClKhrByteAddressableStore struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_byte_addressable_store"`
			ClKhrDeviceUUID struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_device_uuid"`
			ClKhrFp16 struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_fp16"`
			ClKhrGlobalInt32BaseAtomics struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_global_int32_base_atomics"`
			ClKhrGlobalInt32ExtendedAtomics struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_global_int32_extended_atomics"`
			ClKhrIcd struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_icd"`
			ClKhrLocalInt32BaseAtomics struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_local_int32_base_atomics"`
			ClKhrLocalInt32ExtendedAtomics struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_local_int32_extended_atomics"`
			ClIntelCommandQueueFamilies struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_command_queue_families"`
			ClIntelSubgroups struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_subgroups"`
			ClIntelRequiredSubgroupSize struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_required_subgroup_size"`
			ClIntelSubgroupsShort struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_subgroups_short"`
			ClKhrSpir struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_spir"`
			ClIntelAccelerator struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_accelerator"`
			ClIntelDriverDiagnostics struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_driver_diagnostics"`
			ClKhrPriorityHints struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_priority_hints"`
			ClKhrThrottleHints struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_throttle_hints"`
			ClKhrCreateCommandQueue struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_create_command_queue"`
			ClIntelSubgroupsChar struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_subgroups_char"`
			ClIntelSubgroupsLong struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_subgroups_long"`
			ClKhrIlProgram struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_il_program"`
			ClIntelMemForceHostMemory struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_mem_force_host_memory"`
			ClKhrSubgroupExtendedTypes struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_subgroup_extended_types"`
			ClKhrSubgroupNonUniformVote struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_subgroup_non_uniform_vote"`
			ClKhrSubgroupBallot struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_subgroup_ballot"`
			ClKhrSubgroupNonUniformArithmetic struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_subgroup_non_uniform_arithmetic"`
			ClKhrSubgroupShuffle struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_subgroup_shuffle"`
			ClKhrSubgroupShuffleRelative struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_subgroup_shuffle_relative"`
			ClKhrSubgroupClusteredReduce struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_subgroup_clustered_reduce"`
			ClIntelDeviceAttributeQuery struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_device_attribute_query"`
			ClKhrExtendedBitOps struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_extended_bit_ops"`
			ClKhrSuggestedLocalWorkSize struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_suggested_local_work_size"`
			ClIntelSplitWorkGroupBarrier struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_split_work_group_barrier"`
			ClIntelSpirvMediaBlockIo struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_spirv_media_block_io"`
			ClIntelSpirvSubgroups struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_spirv_subgroups"`
			ClKhrSpirvLinkonceOdr struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_spirv_linkonce_odr"`
			ClKhrSpirvNoIntegerWrapDecoration struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_spirv_no_integer_wrap_decoration"`
			ClIntelUnifiedSharedMemory struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_unified_shared_memory"`
			ClKhrMipmapImage struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_mipmap_image"`
			ClKhrMipmapImageWrites struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_mipmap_image_writes"`
			ClExtFloatAtomics struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_ext_float_atomics"`
			ClKhrExternalMemory struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_external_memory"`
			ClIntelPlanarYuv struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_planar_yuv"`
			ClIntelPackedYuv struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_packed_yuv"`
			ClKhrInt64BaseAtomics struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_int64_base_atomics"`
			ClKhrInt64ExtendedAtomics struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_int64_extended_atomics"`
			ClKhrImage2DFromBuffer struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_image2d_from_buffer"`
			ClKhrDepthImages struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_depth_images"`
			ClKhr3DImageWrites struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_3d_image_writes"`
			ClIntelMediaBlockIo struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_media_block_io"`
			ClIntelBfloat16Conversions struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_bfloat16_conversions"`
			ClIntelCreateBufferWithProperties struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_create_buffer_with_properties"`
			ClIntelSubgroupLocalBlockIo struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_subgroup_local_block_io"`
			ClIntelSubgroupMatrixMultiplyAccumulate struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_subgroup_matrix_multiply_accumulate"`
			ClIntelSubgroupSplitMatrixMultiplyAccumulate struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_subgroup_split_matrix_multiply_accumulate"`
			ClKhrIntegerDotProduct struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_integer_dot_product"`
			ClKhrGlSharing struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_gl_sharing"`
			ClKhrGlDepthImages struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_gl_depth_images"`
			ClKhrGlEvent struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_gl_event"`
			ClKhrGlMsaaSharing struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_gl_msaa_sharing"`
			ClIntelVaAPIMediaSharing struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_va_api_media_sharing"`
			ClIntelSharingFormatQuery struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_intel_sharing_format_query"`
			ClKhrPciBusInfo struct {
				Raw     int    `json:"raw"`
				Version string `json:"version"`
			} `json:"cl_khr_pci_bus_info"`
		} `json:"CL_PLATFORM_EXTENSIONS_WITH_VERSION"`
		ClPlatformNumericVersion struct {
			Raw     int    `json:"raw"`
			Version string `json:"version"`
		} `json:"CL_PLATFORM_NUMERIC_VERSION"`
		ClPlatformIcdSuffixKhr                       string   `json:"CL_PLATFORM_ICD_SUFFIX_KHR"`
		ClPlatformHostTimerResolution                int      `json:"CL_PLATFORM_HOST_TIMER_RESOLUTION"`
		ClPlatformExternalMemoryImportHandleTypesKhr []string `json:"CL_PLATFORM_EXTERNAL_MEMORY_IMPORT_HANDLE_TYPES_KHR"`
	} `json:"platforms"`
	Devices []struct {
		Online []struct {
			ClDeviceName           string `json:"CL_DEVICE_NAME"`
			ClDeviceVendor         string `json:"CL_DEVICE_VENDOR"`
			ClDeviceVendorID       int    `json:"CL_DEVICE_VENDOR_ID"`
			ClDeviceVersion        string `json:"CL_DEVICE_VERSION"`
			ClDeviceUUIDKhr        string `json:"CL_DEVICE_UUID_KHR"`
			ClDriverUUIDKhr        string `json:"CL_DRIVER_UUID_KHR"`
			ClDeviceLuidValidKhr   bool   `json:"CL_DEVICE_LUID_VALID_KHR"`
			ClDeviceLuidKhr        string `json:"CL_DEVICE_LUID_KHR"`
			ClDeviceNodeMaskKhr    int    `json:"CL_DEVICE_NODE_MASK_KHR"`
			ClDeviceNumericVersion struct {
				Raw     int    `json:" raw "`
				Version string `json:"version"`
			} `json:"CL_DEVICE_NUMERIC_VERSION"`
			ClDriverVersion            string `json:"CL_DRIVER_VERSION"`
			ClDeviceOpenclCVersion     string `json:"CL_DEVICE_OPENCL_C_VERSION"`
			ClDeviceOpenclCAllVersions struct {
				OpenCLC struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"OpenCL C"`
			} `json:"CL_DEVICE_OPENCL_C_ALL_VERSIONS"`
			ClDeviceOpenclCFeatures struct {
				OpenclCInt64 struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_int64"`
				OpenclC3DImageWrites struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_3d_image_writes"`
				OpenclCImages struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_images"`
				OpenclCReadWriteImages struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_read_write_images"`
				OpenclCAtomicOrderAcqRel struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_atomic_order_acq_rel"`
				OpenclCAtomicOrderSeqCst struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_atomic_order_seq_cst"`
				OpenclCAtomicScopeAllDevices struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_atomic_scope_all_devices"`
				OpenclCAtomicScopeDevice struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_atomic_scope_device"`
				OpenclCGenericAddressSpace struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_generic_address_space"`
				OpenclCProgramScopeGlobalVariables struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_program_scope_global_variables"`
				OpenclCWorkGroupCollectiveFunctions struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_work_group_collective_functions"`
				OpenclCSubgroups struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_subgroups"`
				OpenclCExtFp32GlobalAtomicAdd struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_ext_fp32_global_atomic_add"`
				OpenclCExtFp32LocalAtomicAdd struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_ext_fp32_local_atomic_add"`
				OpenclCExtFp32GlobalAtomicMinMax struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_ext_fp32_global_atomic_min_max"`
				OpenclCExtFp32LocalAtomicMinMax struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_ext_fp32_local_atomic_min_max"`
				OpenclCExtFp16GlobalAtomicLoadStore struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_ext_fp16_global_atomic_load_store"`
				OpenclCExtFp16LocalAtomicLoadStore struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_ext_fp16_local_atomic_load_store"`
				OpenclCExtFp16GlobalAtomicMinMax struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_ext_fp16_global_atomic_min_max"`
				OpenclCExtFp16LocalAtomicMinMax struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_ext_fp16_local_atomic_min_max"`
				OpenclCIntegerDotProductInput4X8Bit struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_integer_dot_product_input_4x8bit"`
				OpenclCIntegerDotProductInput4X8BitPacked struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"__opencl_c_integer_dot_product_input_4x8bit_packed"`
			} `json:"CL_DEVICE_OPENCL_C_FEATURES"`
			ClDeviceLatestConformanceVersionPassed string `json:"CL_DEVICE_LATEST_CONFORMANCE_VERSION_PASSED"`
			ClDeviceType                           struct {
				Raw  int      `json:"raw"`
				Type []string `json:"type"`
			} `json:"CL_DEVICE_TYPE"`
			ClDevicePciBusInfoKhr     string `json:"CL_DEVICE_PCI_BUS_INFO_KHR"`
			ClDeviceProfile           string `json:"CL_DEVICE_PROFILE"`
			ClDeviceAvailable         bool   `json:"CL_DEVICE_AVAILABLE"`
			ClDeviceCompilerAvailable bool   `json:"CL_DEVICE_COMPILER_AVAILABLE"`
			ClDeviceLinkerAvailable   bool   `json:"CL_DEVICE_LINKER_AVAILABLE"`
			ClDeviceMaxComputeUnits   int    `json:"CL_DEVICE_MAX_COMPUTE_UNITS"`
			ClDeviceMaxClockFrequency int    `json:"CL_DEVICE_MAX_CLOCK_FREQUENCY"`
			ClDeviceIPVersionIntel    struct {
				Raw     int    `json:" raw "`
				Version string `json:"version"`
			} `json:"CL_DEVICE_IP_VERSION_INTEL"`
			ClDeviceIDIntel                   int `json:"CL_DEVICE_ID_INTEL"`
			ClDeviceNumSlicesIntel            int `json:"CL_DEVICE_NUM_SLICES_INTEL"`
			ClDeviceNumSubSlicesPerSliceIntel int `json:"CL_DEVICE_NUM_SUB_SLICES_PER_SLICE_INTEL"`
			ClDeviceNumEusPerSubSliceIntel    int `json:"CL_DEVICE_NUM_EUS_PER_SUB_SLICE_INTEL"`
			ClDeviceNumThreadsPerEuIntel      int `json:"CL_DEVICE_NUM_THREADS_PER_EU_INTEL"`
			ClDeviceFeatureCapabilitiesIntel  struct {
				Raw           int      `json:"raw"`
				FeaturesIntel []string `json:"features_intel"`
			} `json:"CL_DEVICE_FEATURE_CAPABILITIES_INTEL"`
			ClDevicePartitionMaxSubDevices  int      `json:"CL_DEVICE_PARTITION_MAX_SUB_DEVICES"`
			ClDevicePartitionProperties     []string `json:"CL_DEVICE_PARTITION_PROPERTIES"`
			ClDevicePartitionAffinityDomain struct {
				Raw    int   `json:"raw"`
				Domain []any `json:"domain"`
			} `json:"CL_DEVICE_PARTITION_AFFINITY_DOMAIN"`
			ClDeviceMaxWorkItemDimensions          int   `json:"CL_DEVICE_MAX_WORK_ITEM_DIMENSIONS"`
			ClDeviceMaxWorkItemSizes               []int `json:"CL_DEVICE_MAX_WORK_ITEM_SIZES"`
			ClDeviceMaxWorkGroupSize               int   `json:"CL_DEVICE_MAX_WORK_GROUP_SIZE"`
			ClDevicePreferredWorkGroupSizeMultiple int   `json:"CL_DEVICE_PREFERRED_WORK_GROUP_SIZE_MULTIPLE"`
			//ClKernelPreferredWorkGroupSizeMultiple int   `json:"CL_KERNEL_PREFERRED_WORK_GROUP_SIZE_MULTIPLE"` // Inside a snap this field is a string with value "<getWGsizes:1954: build program : error -6>"
			ClDeviceMaxNumSubGroups            int   `json:"CL_DEVICE_MAX_NUM_SUB_GROUPS"`
			ClDeviceSubGroupSizesIntel         []int `json:"CL_DEVICE_SUB_GROUP_SIZES_INTEL"`
			ClDevicePreferredVectorWidthChar   int   `json:"CL_DEVICE_PREFERRED_VECTOR_WIDTH_CHAR"`
			ClDeviceNativeVectorWidthChar      int   `json:"CL_DEVICE_NATIVE_VECTOR_WIDTH_CHAR"`
			ClDevicePreferredVectorWidthShort  int   `json:"CL_DEVICE_PREFERRED_VECTOR_WIDTH_SHORT"`
			ClDeviceNativeVectorWidthShort     int   `json:"CL_DEVICE_NATIVE_VECTOR_WIDTH_SHORT"`
			ClDevicePreferredVectorWidthInt    int   `json:"CL_DEVICE_PREFERRED_VECTOR_WIDTH_INT"`
			ClDeviceNativeVectorWidthInt       int   `json:"CL_DEVICE_NATIVE_VECTOR_WIDTH_INT"`
			ClDevicePreferredVectorWidthLong   int   `json:"CL_DEVICE_PREFERRED_VECTOR_WIDTH_LONG"`
			ClDeviceNativeVectorWidthLong      int   `json:"CL_DEVICE_NATIVE_VECTOR_WIDTH_LONG"`
			ClDevicePreferredVectorWidthHalf   int   `json:"CL_DEVICE_PREFERRED_VECTOR_WIDTH_HALF"`
			ClDeviceNativeVectorWidthHalf      int   `json:"CL_DEVICE_NATIVE_VECTOR_WIDTH_HALF"`
			ClDevicePreferredVectorWidthFloat  int   `json:"CL_DEVICE_PREFERRED_VECTOR_WIDTH_FLOAT"`
			ClDeviceNativeVectorWidthFloat     int   `json:"CL_DEVICE_NATIVE_VECTOR_WIDTH_FLOAT"`
			ClDevicePreferredVectorWidthDouble int   `json:"CL_DEVICE_PREFERRED_VECTOR_WIDTH_DOUBLE"`
			ClDeviceNativeVectorWidthDouble    int   `json:"CL_DEVICE_NATIVE_VECTOR_WIDTH_DOUBLE"`
			ClDeviceHalfFpConfig               struct {
				Raw    int      `json:"raw"`
				Config []string `json:"config"`
			} `json:"CL_DEVICE_HALF_FP_CONFIG"`
			ClDeviceSingleFpConfig struct {
				Raw    int      `json:"raw"`
				Config []string `json:"config"`
			} `json:"CL_DEVICE_SINGLE_FP_CONFIG"`
			ClDeviceAddressBits                        int      `json:"CL_DEVICE_ADDRESS_BITS"`
			ClDeviceEndianLittle                       bool     `json:"CL_DEVICE_ENDIAN_LITTLE"`
			ClDeviceExternalMemoryImportHandleTypesKhr []string `json:"CL_DEVICE_EXTERNAL_MEMORY_IMPORT_HANDLE_TYPES_KHR"`
			ClDeviceGlobalMemSize                      uint64   `json:"CL_DEVICE_GLOBAL_MEM_SIZE"`
			ClDeviceErrorCorrectionSupport             bool     `json:"CL_DEVICE_ERROR_CORRECTION_SUPPORT"`
			ClDeviceMaxMemAllocSize                    int64    `json:"CL_DEVICE_MAX_MEM_ALLOC_SIZE"`
			ClDeviceHostUnifiedMemory                  bool     `json:"CL_DEVICE_HOST_UNIFIED_MEMORY"`
			ClDeviceSvmCapabilities                    struct {
				Raw          int      `json:"raw"`
				Capabilities []string `json:"capabilities"`
			} `json:"CL_DEVICE_SVM_CAPABILITIES"`
			ClDeviceHostMemCapabilitiesIntel struct {
				Raw          int      `json:"raw"`
				Capabilities []string `json:"capabilities"`
			} `json:"CL_DEVICE_HOST_MEM_CAPABILITIES_INTEL"`
			ClDeviceDeviceMemCapabilitiesIntel struct {
				Raw          int      `json:"raw"`
				Capabilities []string `json:"capabilities"`
			} `json:"CL_DEVICE_DEVICE_MEM_CAPABILITIES_INTEL"`
			ClDeviceSingleDeviceSharedMemCapabilitiesIntel struct {
				Raw          int      `json:"raw"`
				Capabilities []string `json:"capabilities"`
			} `json:"CL_DEVICE_SINGLE_DEVICE_SHARED_MEM_CAPABILITIES_INTEL"`
			ClDeviceCrossDeviceSharedMemCapabilitiesIntel struct {
				Raw          int      `json:"raw"`
				Capabilities []string `json:"capabilities"`
			} `json:"CL_DEVICE_CROSS_DEVICE_SHARED_MEM_CAPABILITIES_INTEL"`
			ClDeviceSharedSystemMemCapabilitiesIntel struct {
				Raw          int   `json:"raw"`
				Capabilities []any `json:"capabilities"`
			} `json:"CL_DEVICE_SHARED_SYSTEM_MEM_CAPABILITIES_INTEL"`
			ClDeviceMinDataTypeAlignSize             int `json:"CL_DEVICE_MIN_DATA_TYPE_ALIGN_SIZE"`
			ClDeviceMemBaseAddrAlign                 int `json:"CL_DEVICE_MEM_BASE_ADDR_ALIGN"`
			ClDevicePreferredPlatformAtomicAlignment int `json:"CL_DEVICE_PREFERRED_PLATFORM_ATOMIC_ALIGNMENT"`
			ClDevicePreferredGlobalAtomicAlignment   int `json:"CL_DEVICE_PREFERRED_GLOBAL_ATOMIC_ALIGNMENT"`
			ClDevicePreferredLocalAtomicAlignment    int `json:"CL_DEVICE_PREFERRED_LOCAL_ATOMIC_ALIGNMENT"`
			ClDeviceAtomicMemoryCapabilities         struct {
				Raw          int      `json:"raw"`
				Capabilities []string `json:"capabilities"`
			} `json:"CL_DEVICE_ATOMIC_MEMORY_CAPABILITIES"`
			ClDeviceAtomicFenceCapabilities struct {
				Raw          int      `json:"raw"`
				Capabilities []string `json:"capabilities"`
			} `json:"CL_DEVICE_ATOMIC_FENCE_CAPABILITIES"`
			ClDeviceMaxGlobalVariableSize            int    `json:"CL_DEVICE_MAX_GLOBAL_VARIABLE_SIZE"`
			ClDeviceGlobalVariablePreferredTotalSize int64  `json:"CL_DEVICE_GLOBAL_VARIABLE_PREFERRED_TOTAL_SIZE"`
			ClDeviceGlobalMemCacheType               string `json:"CL_DEVICE_GLOBAL_MEM_CACHE_TYPE"`
			ClDeviceGlobalMemCacheSize               int    `json:"CL_DEVICE_GLOBAL_MEM_CACHE_SIZE"`
			ClDeviceGlobalMemCachelineSize           int    `json:"CL_DEVICE_GLOBAL_MEM_CACHELINE_SIZE"`
			ClDeviceImageSupport                     bool   `json:"CL_DEVICE_IMAGE_SUPPORT"`
			ClDeviceMaxSamplers                      int    `json:"CL_DEVICE_MAX_SAMPLERS"`
			ClDeviceImageMaxBufferSize               int    `json:"CL_DEVICE_IMAGE_MAX_BUFFER_SIZE"`
			ClDeviceImageMaxArraySize                int    `json:"CL_DEVICE_IMAGE_MAX_ARRAY_SIZE"`
			ClDeviceImageBaseAddressAlignment        int    `json:"CL_DEVICE_IMAGE_BASE_ADDRESS_ALIGNMENT"`
			ClDeviceImagePitchAlignment              int    `json:"CL_DEVICE_IMAGE_PITCH_ALIGNMENT"`
			ClDeviceImage2DMaxHeight                 int    `json:"CL_DEVICE_IMAGE2D_MAX_HEIGHT"`
			ClDeviceImage2DMaxWidth                  int    `json:"CL_DEVICE_IMAGE2D_MAX_WIDTH"`
			ClDevicePlanarYuvMaxHeightIntel          int    `json:"CL_DEVICE_PLANAR_YUV_MAX_HEIGHT_INTEL"`
			ClDevicePlanarYuvMaxWidthIntel           int    `json:"CL_DEVICE_PLANAR_YUV_MAX_WIDTH_INTEL"`
			ClDeviceImage3DMaxHeight                 int    `json:"CL_DEVICE_IMAGE3D_MAX_HEIGHT"`
			ClDeviceImage3DMaxWidth                  int    `json:"CL_DEVICE_IMAGE3D_MAX_WIDTH"`
			ClDeviceImage3DMaxDepth                  int    `json:"CL_DEVICE_IMAGE3D_MAX_DEPTH"`
			ClDeviceMaxReadImageArgs                 int    `json:"CL_DEVICE_MAX_READ_IMAGE_ARGS"`
			ClDeviceMaxWriteImageArgs                int    `json:"CL_DEVICE_MAX_WRITE_IMAGE_ARGS"`
			ClDeviceMaxReadWriteImageArgs            int    `json:"CL_DEVICE_MAX_READ_WRITE_IMAGE_ARGS"`
			ClDevicePipeSupport                      bool   `json:"CL_DEVICE_PIPE_SUPPORT"`
			ClDeviceMaxPipeArgs                      int    `json:"CL_DEVICE_MAX_PIPE_ARGS"`
			ClDevicePipeMaxActiveReservations        int    `json:"CL_DEVICE_PIPE_MAX_ACTIVE_RESERVATIONS"`
			ClDevicePipeMaxPacketSize                int    `json:"CL_DEVICE_PIPE_MAX_PACKET_SIZE"`
			ClDeviceLocalMemType                     string `json:"CL_DEVICE_LOCAL_MEM_TYPE"`
			ClDeviceLocalMemSize                     int    `json:"CL_DEVICE_LOCAL_MEM_SIZE"`
			ClDeviceMaxConstantArgs                  int    `json:"CL_DEVICE_MAX_CONSTANT_ARGS"`
			ClDeviceMaxConstantBufferSize            int64  `json:"CL_DEVICE_MAX_CONSTANT_BUFFER_SIZE"`
			ClDeviceGenericAddressSpaceSupport       bool   `json:"CL_DEVICE_GENERIC_ADDRESS_SPACE_SUPPORT"`
			ClDeviceMaxParameterSize                 int    `json:"CL_DEVICE_MAX_PARAMETER_SIZE"`
			ClDeviceQueueOnHostProperties            struct {
				Raw       int      `json:"raw"`
				QueueProp []string `json:"queue_prop"`
			} `json:"CL_DEVICE_QUEUE_ON_HOST_PROPERTIES"`
			ClDeviceDeviceEnqueueCapabilities struct {
				Raw          int   `json:"raw"`
				Capabilities []any `json:"capabilities"`
			} `json:"CL_DEVICE_DEVICE_ENQUEUE_CAPABILITIES"`
			ClDeviceQueueOnDeviceProperties struct {
				Raw       int   `json:"raw"`
				QueueProp []any `json:"queue_prop"`
			} `json:"CL_DEVICE_QUEUE_ON_DEVICE_PROPERTIES"`
			ClDeviceQueueOnDevicePreferredSize int `json:"CL_DEVICE_QUEUE_ON_DEVICE_PREFERRED_SIZE"`
			ClDeviceQueueOnDeviceMaxSize       int `json:"CL_DEVICE_QUEUE_ON_DEVICE_MAX_SIZE"`
			ClDeviceMaxOnDeviceQueues          int `json:"CL_DEVICE_MAX_ON_DEVICE_QUEUES"`
			ClDeviceMaxOnDeviceEvents          int `json:"CL_DEVICE_MAX_ON_DEVICE_EVENTS"`
			ClDeviceQueueFamilyPropertiesIntel struct {
				Ccs struct {
					Count       int `json:"count"`
					Proprerties struct {
						Raw        int      `json:"raw"`
						Properties []string `json:"properties"`
					} `json:"proprerties"`
					Capabilities struct {
						Raw          int      `json:"raw"`
						Capabilities []string `json:"capabilities"`
					} `json:"capabilities"`
				} `json:"ccs"`
				Bcs struct {
					Count       int `json:"count"`
					Proprerties struct {
						Raw        int      `json:"raw"`
						Properties []string `json:"properties"`
					} `json:"proprerties"`
					Capabilities struct {
						Raw          int      `json:"raw"`
						Capabilities []string `json:"capabilities"`
					} `json:"capabilities"`
				} `json:"bcs"`
			} `json:"CL_DEVICE_QUEUE_FAMILY_PROPERTIES_INTEL"`
			ClDevicePreferredInteropUserSync bool `json:"CL_DEVICE_PREFERRED_INTEROP_USER_SYNC"`
			ClDeviceProfilingTimerResolution int  `json:"CL_DEVICE_PROFILING_TIMER_RESOLUTION"`
			ClDeviceExecutionCapabilities    struct {
				Raw  int      `json:"raw"`
				Type []string `json:"type"`
			} `json:"CL_DEVICE_EXECUTION_CAPABILITIES"`
			ClDeviceNonUniformWorkGroupSupport          bool   `json:"CL_DEVICE_NON_UNIFORM_WORK_GROUP_SUPPORT"`
			ClDeviceWorkGroupCollectiveFunctionsSupport bool   `json:"CL_DEVICE_WORK_GROUP_COLLECTIVE_FUNCTIONS_SUPPORT"`
			ClDeviceSubGroupIndependentForwardProgress  bool   `json:"CL_DEVICE_SUB_GROUP_INDEPENDENT_FORWARD_PROGRESS"`
			ClDeviceIlVersion                           string `json:"CL_DEVICE_IL_VERSION"`
			ClDeviceIlsWithVersion                      struct {
				SPIRV struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"SPIR-V"`
			} `json:"CL_DEVICE_ILS_WITH_VERSION"`
			ClDeviceSpirVersions              string `json:"CL_DEVICE_SPIR_VERSIONS"`
			ClDevicePrintfBufferSize          int    `json:"CL_DEVICE_PRINTF_BUFFER_SIZE"`
			ClDeviceBuiltInKernels            string `json:"CL_DEVICE_BUILT_IN_KERNELS"`
			ClDeviceBuiltInKernelsWithVersion struct {
			} `json:"CL_DEVICE_BUILT_IN_KERNELS_WITH_VERSION"`
			ClDeviceExtensions            string `json:"CL_DEVICE_EXTENSIONS"`
			ClDeviceExtensionsWithVersion struct {
				ClKhrByteAddressableStore struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_byte_addressable_store"`
				ClKhrDeviceUUID struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_device_uuid"`
				ClKhrFp16 struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_fp16"`
				ClKhrGlobalInt32BaseAtomics struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_global_int32_base_atomics"`
				ClKhrGlobalInt32ExtendedAtomics struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_global_int32_extended_atomics"`
				ClKhrIcd struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_icd"`
				ClKhrLocalInt32BaseAtomics struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_local_int32_base_atomics"`
				ClKhrLocalInt32ExtendedAtomics struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_local_int32_extended_atomics"`
				ClIntelCommandQueueFamilies struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_command_queue_families"`
				ClIntelSubgroups struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_subgroups"`
				ClIntelRequiredSubgroupSize struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_required_subgroup_size"`
				ClIntelSubgroupsShort struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_subgroups_short"`
				ClKhrSpir struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_spir"`
				ClIntelAccelerator struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_accelerator"`
				ClIntelDriverDiagnostics struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_driver_diagnostics"`
				ClKhrPriorityHints struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_priority_hints"`
				ClKhrThrottleHints struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_throttle_hints"`
				ClKhrCreateCommandQueue struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_create_command_queue"`
				ClIntelSubgroupsChar struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_subgroups_char"`
				ClIntelSubgroupsLong struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_subgroups_long"`
				ClKhrIlProgram struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_il_program"`
				ClIntelMemForceHostMemory struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_mem_force_host_memory"`
				ClKhrSubgroupExtendedTypes struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_subgroup_extended_types"`
				ClKhrSubgroupNonUniformVote struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_subgroup_non_uniform_vote"`
				ClKhrSubgroupBallot struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_subgroup_ballot"`
				ClKhrSubgroupNonUniformArithmetic struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_subgroup_non_uniform_arithmetic"`
				ClKhrSubgroupShuffle struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_subgroup_shuffle"`
				ClKhrSubgroupShuffleRelative struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_subgroup_shuffle_relative"`
				ClKhrSubgroupClusteredReduce struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_subgroup_clustered_reduce"`
				ClIntelDeviceAttributeQuery struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_device_attribute_query"`
				ClKhrExtendedBitOps struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_extended_bit_ops"`
				ClKhrSuggestedLocalWorkSize struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_suggested_local_work_size"`
				ClIntelSplitWorkGroupBarrier struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_split_work_group_barrier"`
				ClIntelSpirvMediaBlockIo struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_spirv_media_block_io"`
				ClIntelSpirvSubgroups struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_spirv_subgroups"`
				ClKhrSpirvLinkonceOdr struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_spirv_linkonce_odr"`
				ClKhrSpirvNoIntegerWrapDecoration struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_spirv_no_integer_wrap_decoration"`
				ClIntelUnifiedSharedMemory struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_unified_shared_memory"`
				ClKhrMipmapImage struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_mipmap_image"`
				ClKhrMipmapImageWrites struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_mipmap_image_writes"`
				ClExtFloatAtomics struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_ext_float_atomics"`
				ClKhrExternalMemory struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_external_memory"`
				ClIntelPlanarYuv struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_planar_yuv"`
				ClIntelPackedYuv struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_packed_yuv"`
				ClKhrInt64BaseAtomics struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_int64_base_atomics"`
				ClKhrInt64ExtendedAtomics struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_int64_extended_atomics"`
				ClKhrImage2DFromBuffer struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_image2d_from_buffer"`
				ClKhrDepthImages struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_depth_images"`
				ClKhr3DImageWrites struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_3d_image_writes"`
				ClIntelMediaBlockIo struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_media_block_io"`
				ClIntelBfloat16Conversions struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_bfloat16_conversions"`
				ClIntelCreateBufferWithProperties struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_create_buffer_with_properties"`
				ClIntelSubgroupLocalBlockIo struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_subgroup_local_block_io"`
				ClIntelSubgroupMatrixMultiplyAccumulate struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_subgroup_matrix_multiply_accumulate"`
				ClIntelSubgroupSplitMatrixMultiplyAccumulate struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_subgroup_split_matrix_multiply_accumulate"`
				ClKhrIntegerDotProduct struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_integer_dot_product"`
				ClKhrGlSharing struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_gl_sharing"`
				ClKhrGlDepthImages struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_gl_depth_images"`
				ClKhrGlEvent struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_gl_event"`
				ClKhrGlMsaaSharing struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_gl_msaa_sharing"`
				ClIntelVaAPIMediaSharing struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_va_api_media_sharing"`
				ClIntelSharingFormatQuery struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_intel_sharing_format_query"`
				ClKhrPciBusInfo struct {
					Raw     int    `json:"raw"`
					Version string `json:"version"`
				} `json:"cl_khr_pci_bus_info"`
			} `json:"CL_DEVICE_EXTENSIONS_WITH_VERSION"`
		} `json:"online"`
	} `json:"devices"`
	IcdLoader struct {
		ClIcdlName       string `json:"CL_ICDL_NAME"`
		ClIcdlVendor     string `json:"CL_ICDL_VENDOR"`
		ClIcdlVersion    string `json:"CL_ICDL_VERSION"`
		ClIcdlOclVersion string `json:"CL_ICDL_OCL_VERSION"`
		DetectedVersion  string `json:"_detected_version"`
	} `json:"icd_loader"`
}
