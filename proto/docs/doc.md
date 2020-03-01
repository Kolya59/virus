# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [common.proto](#common.proto)
    - [HealthCheckReq](#pb.HealthCheckReq)
    - [HealthCheckRes](#pb.HealthCheckRes)
  
    - [HealthCheckRes.ServingStatus](#pb.HealthCheckRes.ServingStatus)
  
  
  

- [dispatcherFunction.proto](#dispatcherFunction.proto)
    - [GetTargetReq](#pb.GetTargetReq)
    - [GetTargetRes](#pb.GetTargetRes)
  
  
  
    - [FunctionDispatcher](#pb.FunctionDispatcher)
  

- [dispatcherServer.proto](#dispatcherServer.proto)
    - [GetNextServerReq](#pb.GetNextServerReq)
    - [GetNextServerRes](#pb.GetNextServerRes)
    - [RegisterReq](#pb.RegisterReq)
    - [RegisterRes](#pb.RegisterRes)
  
    - [GetNextServerRes.RequestStatus](#pb.GetNextServerRes.RequestStatus)
    - [RegisterRes.RequestStatus](#pb.RegisterRes.RequestStatus)
  
  
    - [ServerDispatcher](#pb.ServerDispatcher)
  

- [workerServer.proto](#workerServer.proto)
    - [Address](#pb.Address)
    - [Iface](#pb.Iface)
    - [Machine](#pb.Machine)
    - [Port](#pb.Port)
    - [SaveMachineReq](#pb.SaveMachineReq)
    - [SaveMachineRes](#pb.SaveMachineRes)
  
    - [SaveMachineRes.RequestStatus](#pb.SaveMachineRes.RequestStatus)
  
  
    - [WorkerSaver](#pb.WorkerSaver)
  

- [Scalar Value Types](#scalar-value-types)



<a name="common.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## common.proto



<a name="pb.HealthCheckReq"></a>

### HealthCheckReq







<a name="pb.HealthCheckRes"></a>

### HealthCheckRes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [HealthCheckRes.ServingStatus](#pb.HealthCheckRes.ServingStatus) |  |  |





 


<a name="pb.HealthCheckRes.ServingStatus"></a>

### HealthCheckRes.ServingStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| SERVING | 1 |  |
| NOT_SERVING | 2 |  |


 

 

 



<a name="dispatcherFunction.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## dispatcherFunction.proto



<a name="pb.GetTargetReq"></a>

### GetTargetReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [bytes](#bytes) |  |  |






<a name="pb.GetTargetRes"></a>

### GetTargetRes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ip | [string](#string) |  |  |
| certificate | [bytes](#bytes) |  |  |





 

 

 


<a name="pb.FunctionDispatcher"></a>

### FunctionDispatcher


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Check | [HealthCheckReq](#pb.HealthCheckReq) | [HealthCheckRes](#pb.HealthCheckRes) |  |
| GetTarget | [GetTargetReq](#pb.GetTargetReq) | [GetTargetRes](#pb.GetTargetRes) |  |

 



<a name="dispatcherServer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## dispatcherServer.proto



<a name="pb.GetNextServerReq"></a>

### GetNextServerReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ip | [string](#string) |  |  |
| certificate | [bytes](#bytes) |  |  |






<a name="pb.GetNextServerRes"></a>

### GetNextServerRes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [GetNextServerRes.RequestStatus](#pb.GetNextServerRes.RequestStatus) |  |  |






<a name="pb.RegisterReq"></a>

### RegisterReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| certificate | [bytes](#bytes) |  |  |






<a name="pb.RegisterRes"></a>

### RegisterRes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [RegisterRes.RequestStatus](#pb.RegisterRes.RequestStatus) |  |  |





 


<a name="pb.GetNextServerRes.RequestStatus"></a>

### GetNextServerRes.RequestStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| ACCEPTED | 1 |  |
| NOT_ACCEPTED | 2 |  |



<a name="pb.RegisterRes.RequestStatus"></a>

### RegisterRes.RequestStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| REGISTERED | 1 |  |
| NOT_REGISTERED | 2 |  |


 

 


<a name="pb.ServerDispatcher"></a>

### ServerDispatcher


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Check | [HealthCheckReq](#pb.HealthCheckReq) | [HealthCheckRes](#pb.HealthCheckRes) |  |
| Register | [RegisterReq](#pb.RegisterReq) | [RegisterRes](#pb.RegisterRes) |  |
| GetNextServer | [GetNextServerReq](#pb.GetNextServerReq) | [GetNextServerRes](#pb.GetNextServerRes) |  |

 



<a name="workerServer.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## workerServer.proto



<a name="pb.Address"></a>

### Address



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ip | [string](#string) |  |  |
| ports | [Port](#pb.Port) | repeated |  |






<a name="pb.Iface"></a>

### Iface



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| address | [Address](#pb.Address) | repeated |  |
| error | [string](#string) |  |  |






<a name="pb.Machine"></a>

### Machine



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uuid | [string](#string) |  |  |
| ifaces | [Iface](#pb.Iface) | repeated |  |
| error | [string](#string) |  |  |






<a name="pb.Port"></a>

### Port



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| port | [uint32](#uint32) |  |  |
| service | [string](#string) |  |  |






<a name="pb.SaveMachineReq"></a>

### SaveMachineReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| machine | [Machine](#pb.Machine) |  |  |






<a name="pb.SaveMachineRes"></a>

### SaveMachineRes



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [SaveMachineRes.RequestStatus](#pb.SaveMachineRes.RequestStatus) |  |  |





 


<a name="pb.SaveMachineRes.RequestStatus"></a>

### SaveMachineRes.RequestStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 |  |
| ACCEPTED | 1 |  |
| NOT_ACCEPTED | 2 |  |


 

 


<a name="pb.WorkerSaver"></a>

### WorkerSaver


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Check | [HealthCheckReq](#pb.HealthCheckReq) | [HealthCheckRes](#pb.HealthCheckRes) |  |
| SaveMachine | [SaveMachineReq](#pb.SaveMachineReq) | [SaveMachineRes](#pb.SaveMachineRes) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

