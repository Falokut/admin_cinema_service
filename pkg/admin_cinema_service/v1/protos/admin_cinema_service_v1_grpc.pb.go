// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: admin_cinema_service_v1.proto

package protos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CinemaServiceV1Client is the client API for CinemaServiceV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CinemaServiceV1Client interface {
	// Returns all cities where there are cinemas.
	GetCinemasCities(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Cities, error)
	// Returns cinemas in the city.
	GetCinemasInCity(ctx context.Context, in *GetCinemasInCityRequest, opts ...grpc.CallOption) (*Cinemas, error)
	// Returns all movies that are in the cinema screenings in a particular cinema.
	GetMoviesScreenings(ctx context.Context, in *GetMoviesScreeningsRequest, opts ...grpc.CallOption) (*PreviewScreenings, error)
	// Returns all screenings for a movie in a specific cinema.
	GetScreenings(ctx context.Context, in *GetScreeningsRequest, opts ...grpc.CallOption) (*Screenings, error)
	// Returns the configuration of the hall.
	GetHallConfiguration(ctx context.Context, in *GetHallConfigurationRequest, opts ...grpc.CallOption) (*HallConfiguration, error)
	// Create screenings, returns created screenings ids.
	CreateScreenings(ctx context.Context, in *CreateScreeningsRequest, opts ...grpc.CallOption) (*CreateScreeningsResponse, error)
	// Create cinema, returns created cinema id.
	CreateCinema(ctx context.Context, in *CreateCinemaRequest, opts ...grpc.CallOption) (*CreateCinemaResponse, error)
	// Delete cinema with the specified id.
	DeleteCinema(ctx context.Context, in *DeleteCinemaRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Update cinema with the specified id.
	UpdateCinema(ctx context.Context, in *UpdateCinemaRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Create city, returns created city id.
	CreateCity(ctx context.Context, in *CreateCityRequest, opts ...grpc.CallOption) (*CreateCityResponse, error)
	// Update city, returns created city id.
	UpdateCity(ctx context.Context, in *UpdateCityRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Delete city with the specified id.
	DeleteCity(ctx context.Context, in *DeleteCityRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Returns halls info without configuration
	GetHalls(ctx context.Context, in *GetHallsRequest, opts ...grpc.CallOption) (*Halls, error)
	// Create hall and configuration for it, configuration mustn't be empty
	CreateHall(ctx context.Context, in *CreateHallRequest, opts ...grpc.CallOption) (*CreateHallResponse, error)
	// Delete hall and hall configuration for specified hall id.
	DeleteHall(ctx context.Context, in *DeleteHallRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Update hall and hall configuration for specified hall id.
	UpdateHall(ctx context.Context, in *UpdateHallRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Create hall type, returns created hall type id.
	CreateHallType(ctx context.Context, in *CreateHallTypeRequest, opts ...grpc.CallOption) (*CreateHallTypeResponse, error)
	// Delete hall type with the specified hall type id.
	DeleteHallType(ctx context.Context, in *DeleteHallTypeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Update hall type(set new name) with the specified hall type id.
	UpdateHallType(ctx context.Context, in *UpdateHallTypeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Create screening type, returns created screening type id.
	CreateScreeningType(ctx context.Context, in *CreateScreeningTypeRequest, opts ...grpc.CallOption) (*CreateScreeningTypeResponse, error)
	// Delete screening type with the specified screening type id.
	DeleteScreeningType(ctx context.Context, in *DeleteScreeningTypeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Update screening type(set new name) with the specified screening type id.
	UpdateScreeningType(ctx context.Context, in *UpdateScreeningTypeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type cinemaServiceV1Client struct {
	cc grpc.ClientConnInterface
}

func NewCinemaServiceV1Client(cc grpc.ClientConnInterface) CinemaServiceV1Client {
	return &cinemaServiceV1Client{cc}
}

func (c *cinemaServiceV1Client) GetCinemasCities(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Cities, error) {
	out := new(Cities)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/GetCinemasCities", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) GetCinemasInCity(ctx context.Context, in *GetCinemasInCityRequest, opts ...grpc.CallOption) (*Cinemas, error) {
	out := new(Cinemas)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/GetCinemasInCity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) GetMoviesScreenings(ctx context.Context, in *GetMoviesScreeningsRequest, opts ...grpc.CallOption) (*PreviewScreenings, error) {
	out := new(PreviewScreenings)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/GetMoviesScreenings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) GetScreenings(ctx context.Context, in *GetScreeningsRequest, opts ...grpc.CallOption) (*Screenings, error) {
	out := new(Screenings)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/GetScreenings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) GetHallConfiguration(ctx context.Context, in *GetHallConfigurationRequest, opts ...grpc.CallOption) (*HallConfiguration, error) {
	out := new(HallConfiguration)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/GetHallConfiguration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) CreateScreenings(ctx context.Context, in *CreateScreeningsRequest, opts ...grpc.CallOption) (*CreateScreeningsResponse, error) {
	out := new(CreateScreeningsResponse)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/CreateScreenings", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) CreateCinema(ctx context.Context, in *CreateCinemaRequest, opts ...grpc.CallOption) (*CreateCinemaResponse, error) {
	out := new(CreateCinemaResponse)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/CreateCinema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) DeleteCinema(ctx context.Context, in *DeleteCinemaRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/DeleteCinema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) UpdateCinema(ctx context.Context, in *UpdateCinemaRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/UpdateCinema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) CreateCity(ctx context.Context, in *CreateCityRequest, opts ...grpc.CallOption) (*CreateCityResponse, error) {
	out := new(CreateCityResponse)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/CreateCity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) UpdateCity(ctx context.Context, in *UpdateCityRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/UpdateCity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) DeleteCity(ctx context.Context, in *DeleteCityRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/DeleteCity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) GetHalls(ctx context.Context, in *GetHallsRequest, opts ...grpc.CallOption) (*Halls, error) {
	out := new(Halls)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/GetHalls", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) CreateHall(ctx context.Context, in *CreateHallRequest, opts ...grpc.CallOption) (*CreateHallResponse, error) {
	out := new(CreateHallResponse)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/CreateHall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) DeleteHall(ctx context.Context, in *DeleteHallRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/DeleteHall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) UpdateHall(ctx context.Context, in *UpdateHallRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/UpdateHall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) CreateHallType(ctx context.Context, in *CreateHallTypeRequest, opts ...grpc.CallOption) (*CreateHallTypeResponse, error) {
	out := new(CreateHallTypeResponse)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/CreateHallType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) DeleteHallType(ctx context.Context, in *DeleteHallTypeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/DeleteHallType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) UpdateHallType(ctx context.Context, in *UpdateHallTypeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/UpdateHallType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) CreateScreeningType(ctx context.Context, in *CreateScreeningTypeRequest, opts ...grpc.CallOption) (*CreateScreeningTypeResponse, error) {
	out := new(CreateScreeningTypeResponse)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/CreateScreeningType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) DeleteScreeningType(ctx context.Context, in *DeleteScreeningTypeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/DeleteScreeningType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cinemaServiceV1Client) UpdateScreeningType(ctx context.Context, in *UpdateScreeningTypeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/admin_cinema_service.cinemaServiceV1/UpdateScreeningType", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CinemaServiceV1Server is the server API for CinemaServiceV1 service.
// All implementations must embed UnimplementedCinemaServiceV1Server
// for forward compatibility
type CinemaServiceV1Server interface {
	// Returns all cities where there are cinemas.
	GetCinemasCities(context.Context, *emptypb.Empty) (*Cities, error)
	// Returns cinemas in the city.
	GetCinemasInCity(context.Context, *GetCinemasInCityRequest) (*Cinemas, error)
	// Returns all movies that are in the cinema screenings in a particular cinema.
	GetMoviesScreenings(context.Context, *GetMoviesScreeningsRequest) (*PreviewScreenings, error)
	// Returns all screenings for a movie in a specific cinema.
	GetScreenings(context.Context, *GetScreeningsRequest) (*Screenings, error)
	// Returns the configuration of the hall.
	GetHallConfiguration(context.Context, *GetHallConfigurationRequest) (*HallConfiguration, error)
	// Create screenings, returns created screenings ids.
	CreateScreenings(context.Context, *CreateScreeningsRequest) (*CreateScreeningsResponse, error)
	// Create cinema, returns created cinema id.
	CreateCinema(context.Context, *CreateCinemaRequest) (*CreateCinemaResponse, error)
	// Delete cinema with the specified id.
	DeleteCinema(context.Context, *DeleteCinemaRequest) (*emptypb.Empty, error)
	// Update cinema with the specified id.
	UpdateCinema(context.Context, *UpdateCinemaRequest) (*emptypb.Empty, error)
	// Create city, returns created city id.
	CreateCity(context.Context, *CreateCityRequest) (*CreateCityResponse, error)
	// Update city, returns created city id.
	UpdateCity(context.Context, *UpdateCityRequest) (*emptypb.Empty, error)
	// Delete city with the specified id.
	DeleteCity(context.Context, *DeleteCityRequest) (*emptypb.Empty, error)
	// Returns halls info without configuration
	GetHalls(context.Context, *GetHallsRequest) (*Halls, error)
	// Create hall and configuration for it, configuration mustn't be empty
	CreateHall(context.Context, *CreateHallRequest) (*CreateHallResponse, error)
	// Delete hall and hall configuration for specified hall id.
	DeleteHall(context.Context, *DeleteHallRequest) (*emptypb.Empty, error)
	// Update hall and hall configuration for specified hall id.
	UpdateHall(context.Context, *UpdateHallRequest) (*emptypb.Empty, error)
	// Create hall type, returns created hall type id.
	CreateHallType(context.Context, *CreateHallTypeRequest) (*CreateHallTypeResponse, error)
	// Delete hall type with the specified hall type id.
	DeleteHallType(context.Context, *DeleteHallTypeRequest) (*emptypb.Empty, error)
	// Update hall type(set new name) with the specified hall type id.
	UpdateHallType(context.Context, *UpdateHallTypeRequest) (*emptypb.Empty, error)
	// Create screening type, returns created screening type id.
	CreateScreeningType(context.Context, *CreateScreeningTypeRequest) (*CreateScreeningTypeResponse, error)
	// Delete screening type with the specified screening type id.
	DeleteScreeningType(context.Context, *DeleteScreeningTypeRequest) (*emptypb.Empty, error)
	// Update screening type(set new name) with the specified screening type id.
	UpdateScreeningType(context.Context, *UpdateScreeningTypeRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedCinemaServiceV1Server()
}

// UnimplementedCinemaServiceV1Server must be embedded to have forward compatible implementations.
type UnimplementedCinemaServiceV1Server struct {
}

func (UnimplementedCinemaServiceV1Server) GetCinemasCities(context.Context, *emptypb.Empty) (*Cities, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCinemasCities not implemented")
}
func (UnimplementedCinemaServiceV1Server) GetCinemasInCity(context.Context, *GetCinemasInCityRequest) (*Cinemas, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCinemasInCity not implemented")
}
func (UnimplementedCinemaServiceV1Server) GetMoviesScreenings(context.Context, *GetMoviesScreeningsRequest) (*PreviewScreenings, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMoviesScreenings not implemented")
}
func (UnimplementedCinemaServiceV1Server) GetScreenings(context.Context, *GetScreeningsRequest) (*Screenings, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetScreenings not implemented")
}
func (UnimplementedCinemaServiceV1Server) GetHallConfiguration(context.Context, *GetHallConfigurationRequest) (*HallConfiguration, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHallConfiguration not implemented")
}
func (UnimplementedCinemaServiceV1Server) CreateScreenings(context.Context, *CreateScreeningsRequest) (*CreateScreeningsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateScreenings not implemented")
}
func (UnimplementedCinemaServiceV1Server) CreateCinema(context.Context, *CreateCinemaRequest) (*CreateCinemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCinema not implemented")
}
func (UnimplementedCinemaServiceV1Server) DeleteCinema(context.Context, *DeleteCinemaRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCinema not implemented")
}
func (UnimplementedCinemaServiceV1Server) UpdateCinema(context.Context, *UpdateCinemaRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCinema not implemented")
}
func (UnimplementedCinemaServiceV1Server) CreateCity(context.Context, *CreateCityRequest) (*CreateCityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCity not implemented")
}
func (UnimplementedCinemaServiceV1Server) UpdateCity(context.Context, *UpdateCityRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCity not implemented")
}
func (UnimplementedCinemaServiceV1Server) DeleteCity(context.Context, *DeleteCityRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCity not implemented")
}
func (UnimplementedCinemaServiceV1Server) GetHalls(context.Context, *GetHallsRequest) (*Halls, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHalls not implemented")
}
func (UnimplementedCinemaServiceV1Server) CreateHall(context.Context, *CreateHallRequest) (*CreateHallResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateHall not implemented")
}
func (UnimplementedCinemaServiceV1Server) DeleteHall(context.Context, *DeleteHallRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteHall not implemented")
}
func (UnimplementedCinemaServiceV1Server) UpdateHall(context.Context, *UpdateHallRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateHall not implemented")
}
func (UnimplementedCinemaServiceV1Server) CreateHallType(context.Context, *CreateHallTypeRequest) (*CreateHallTypeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateHallType not implemented")
}
func (UnimplementedCinemaServiceV1Server) DeleteHallType(context.Context, *DeleteHallTypeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteHallType not implemented")
}
func (UnimplementedCinemaServiceV1Server) UpdateHallType(context.Context, *UpdateHallTypeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateHallType not implemented")
}
func (UnimplementedCinemaServiceV1Server) CreateScreeningType(context.Context, *CreateScreeningTypeRequest) (*CreateScreeningTypeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateScreeningType not implemented")
}
func (UnimplementedCinemaServiceV1Server) DeleteScreeningType(context.Context, *DeleteScreeningTypeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteScreeningType not implemented")
}
func (UnimplementedCinemaServiceV1Server) UpdateScreeningType(context.Context, *UpdateScreeningTypeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateScreeningType not implemented")
}
func (UnimplementedCinemaServiceV1Server) mustEmbedUnimplementedCinemaServiceV1Server() {}

// UnsafeCinemaServiceV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CinemaServiceV1Server will
// result in compilation errors.
type UnsafeCinemaServiceV1Server interface {
	mustEmbedUnimplementedCinemaServiceV1Server()
}

func RegisterCinemaServiceV1Server(s grpc.ServiceRegistrar, srv CinemaServiceV1Server) {
	s.RegisterService(&CinemaServiceV1_ServiceDesc, srv)
}

func _CinemaServiceV1_GetCinemasCities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).GetCinemasCities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/GetCinemasCities",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).GetCinemasCities(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_GetCinemasInCity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCinemasInCityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).GetCinemasInCity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/GetCinemasInCity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).GetCinemasInCity(ctx, req.(*GetCinemasInCityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_GetMoviesScreenings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMoviesScreeningsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).GetMoviesScreenings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/GetMoviesScreenings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).GetMoviesScreenings(ctx, req.(*GetMoviesScreeningsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_GetScreenings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetScreeningsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).GetScreenings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/GetScreenings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).GetScreenings(ctx, req.(*GetScreeningsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_GetHallConfiguration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetHallConfigurationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).GetHallConfiguration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/GetHallConfiguration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).GetHallConfiguration(ctx, req.(*GetHallConfigurationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_CreateScreenings_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateScreeningsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).CreateScreenings(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/CreateScreenings",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).CreateScreenings(ctx, req.(*CreateScreeningsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_CreateCinema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCinemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).CreateCinema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/CreateCinema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).CreateCinema(ctx, req.(*CreateCinemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_DeleteCinema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCinemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).DeleteCinema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/DeleteCinema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).DeleteCinema(ctx, req.(*DeleteCinemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_UpdateCinema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCinemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).UpdateCinema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/UpdateCinema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).UpdateCinema(ctx, req.(*UpdateCinemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_CreateCity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).CreateCity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/CreateCity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).CreateCity(ctx, req.(*CreateCityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_UpdateCity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).UpdateCity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/UpdateCity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).UpdateCity(ctx, req.(*UpdateCityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_DeleteCity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).DeleteCity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/DeleteCity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).DeleteCity(ctx, req.(*DeleteCityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_GetHalls_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetHallsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).GetHalls(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/GetHalls",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).GetHalls(ctx, req.(*GetHallsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_CreateHall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateHallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).CreateHall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/CreateHall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).CreateHall(ctx, req.(*CreateHallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_DeleteHall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteHallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).DeleteHall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/DeleteHall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).DeleteHall(ctx, req.(*DeleteHallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_UpdateHall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateHallRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).UpdateHall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/UpdateHall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).UpdateHall(ctx, req.(*UpdateHallRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_CreateHallType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateHallTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).CreateHallType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/CreateHallType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).CreateHallType(ctx, req.(*CreateHallTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_DeleteHallType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteHallTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).DeleteHallType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/DeleteHallType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).DeleteHallType(ctx, req.(*DeleteHallTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_UpdateHallType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateHallTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).UpdateHallType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/UpdateHallType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).UpdateHallType(ctx, req.(*UpdateHallTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_CreateScreeningType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateScreeningTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).CreateScreeningType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/CreateScreeningType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).CreateScreeningType(ctx, req.(*CreateScreeningTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_DeleteScreeningType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteScreeningTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).DeleteScreeningType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/DeleteScreeningType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).DeleteScreeningType(ctx, req.(*DeleteScreeningTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CinemaServiceV1_UpdateScreeningType_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateScreeningTypeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CinemaServiceV1Server).UpdateScreeningType(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/admin_cinema_service.cinemaServiceV1/UpdateScreeningType",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CinemaServiceV1Server).UpdateScreeningType(ctx, req.(*UpdateScreeningTypeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CinemaServiceV1_ServiceDesc is the grpc.ServiceDesc for CinemaServiceV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CinemaServiceV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "admin_cinema_service.cinemaServiceV1",
	HandlerType: (*CinemaServiceV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCinemasCities",
			Handler:    _CinemaServiceV1_GetCinemasCities_Handler,
		},
		{
			MethodName: "GetCinemasInCity",
			Handler:    _CinemaServiceV1_GetCinemasInCity_Handler,
		},
		{
			MethodName: "GetMoviesScreenings",
			Handler:    _CinemaServiceV1_GetMoviesScreenings_Handler,
		},
		{
			MethodName: "GetScreenings",
			Handler:    _CinemaServiceV1_GetScreenings_Handler,
		},
		{
			MethodName: "GetHallConfiguration",
			Handler:    _CinemaServiceV1_GetHallConfiguration_Handler,
		},
		{
			MethodName: "CreateScreenings",
			Handler:    _CinemaServiceV1_CreateScreenings_Handler,
		},
		{
			MethodName: "CreateCinema",
			Handler:    _CinemaServiceV1_CreateCinema_Handler,
		},
		{
			MethodName: "DeleteCinema",
			Handler:    _CinemaServiceV1_DeleteCinema_Handler,
		},
		{
			MethodName: "UpdateCinema",
			Handler:    _CinemaServiceV1_UpdateCinema_Handler,
		},
		{
			MethodName: "CreateCity",
			Handler:    _CinemaServiceV1_CreateCity_Handler,
		},
		{
			MethodName: "UpdateCity",
			Handler:    _CinemaServiceV1_UpdateCity_Handler,
		},
		{
			MethodName: "DeleteCity",
			Handler:    _CinemaServiceV1_DeleteCity_Handler,
		},
		{
			MethodName: "GetHalls",
			Handler:    _CinemaServiceV1_GetHalls_Handler,
		},
		{
			MethodName: "CreateHall",
			Handler:    _CinemaServiceV1_CreateHall_Handler,
		},
		{
			MethodName: "DeleteHall",
			Handler:    _CinemaServiceV1_DeleteHall_Handler,
		},
		{
			MethodName: "UpdateHall",
			Handler:    _CinemaServiceV1_UpdateHall_Handler,
		},
		{
			MethodName: "CreateHallType",
			Handler:    _CinemaServiceV1_CreateHallType_Handler,
		},
		{
			MethodName: "DeleteHallType",
			Handler:    _CinemaServiceV1_DeleteHallType_Handler,
		},
		{
			MethodName: "UpdateHallType",
			Handler:    _CinemaServiceV1_UpdateHallType_Handler,
		},
		{
			MethodName: "CreateScreeningType",
			Handler:    _CinemaServiceV1_CreateScreeningType_Handler,
		},
		{
			MethodName: "DeleteScreeningType",
			Handler:    _CinemaServiceV1_DeleteScreeningType_Handler,
		},
		{
			MethodName: "UpdateScreeningType",
			Handler:    _CinemaServiceV1_UpdateScreeningType_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "admin_cinema_service_v1.proto",
}