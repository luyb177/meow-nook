// user_ext.pb.go — hand-written extension messages for VerifyCode and Register RPCs.
// These supplement the protoc-generated user.pb.go without touching its rawDesc.

package v1

// VerifyCodeReq is the request for the VerifyCode RPC.
type VerifyCodeReq struct {
	Target  string        `protobuf:"bytes,1,opt,name=target,proto3" json:"target,omitempty"`
	Channel VerifyChannel `protobuf:"varint,2,opt,name=channel,proto3,enum=user.VerifyChannel" json:"channel,omitempty"`
	Purpose VerifyPurpose `protobuf:"varint,3,opt,name=purpose,proto3,enum=user.VerifyPurpose" json:"purpose,omitempty"`
	Code    string        `protobuf:"bytes,4,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *VerifyCodeReq) Reset()         {}
func (x *VerifyCodeReq) String() string  { return x.Target }
func (x *VerifyCodeReq) ProtoMessage()  {}

func (x *VerifyCodeReq) GetTarget() string        { return x.Target }
func (x *VerifyCodeReq) GetChannel() VerifyChannel { return x.Channel }
func (x *VerifyCodeReq) GetPurpose() VerifyPurpose { return x.Purpose }
func (x *VerifyCodeReq) GetCode() string           { return x.Code }

// RegisterReq is the request for the Register RPC.
type RegisterReq struct {
	Email    string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *RegisterReq) Reset()         {}
func (x *RegisterReq) String() string  { return x.Email }
func (x *RegisterReq) ProtoMessage()  {}

func (x *RegisterReq) GetEmail() string    { return x.Email }
func (x *RegisterReq) GetPassword() string { return x.Password }

// RegisterResp is the response for the Register RPC.
type RegisterResp struct {
	UserId int64  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *RegisterResp) Reset()        {}
func (x *RegisterResp) String() string { return "" }
func (x *RegisterResp) ProtoMessage() {}

func (x *RegisterResp) GetUserId() int64  { return x.UserId }
func (x *RegisterResp) GetToken() string  { return x.Token }

// LoginReq is the request for the Login RPC.
type LoginReq struct {
	Email    string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *LoginReq) Reset()        {}
func (x *LoginReq) String() string { return x.Email }
func (x *LoginReq) ProtoMessage() {}

func (x *LoginReq) GetEmail() string    { return x.Email }
func (x *LoginReq) GetPassword() string { return x.Password }

// LoginResp is the response for the Login RPC.
type LoginResp struct {
	Token  string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	UserId int64  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *LoginResp) Reset()        {}
func (x *LoginResp) String() string { return "" }
func (x *LoginResp) ProtoMessage() {}

func (x *LoginResp) GetToken() string  { return x.Token }
func (x *LoginResp) GetUserId() int64  { return x.UserId }
