* grpc验证器和pb.go文件    
  protoc -I . --govalidators_out=. --go_out=plugins=grpc:. simple.proto
* grpc-gateway:    
protoc --grpc-gateway_out=logtostderr=true:./ ./simple.proto

* oenapiv2:   
 protoc  --openapiv2_out . --openapiv2_opt logtostderr=true  simple.proto
 
 
 具体参考连接: 
 https://www.cnblogs.com/FireworksEasyCool/p/12782137.html

