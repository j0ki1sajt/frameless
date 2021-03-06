package frameless

import (
	"context"
	"io"
	"testing"
)

// Entity encapsulate the most general and high-level rules of the application.
// 	"An entity can be an object with methods, or it can be a set of data structures and functions"
// 	Robert Martin
//
// In enterprise environment, this or the specification of this object can be shared between applications.
// If you don’t have an enterprise, and are just writing a single application, then these entities are the business objects of the application.
// They encapsulate the most general and high-level rules.
// Entity scope must be free from anything related to other software layers implementation knowledge such as SQL or HTTP request objects.
//
// They are the least likely to change when something external changes.
// For example, you would not expect these objects to be affected by a change to page navigation, or security.
// No operational change to any particular application should affect the entity layer.
//
// By convention these structures should be placed on the top folder level of the project
type Entity = interface{}

// Controller defines how a framework independent controller should look
//  the core concept that the controller implements the business rules
// 	and the channels like CLI or HTTP only provide data to it in a form of application specific business entities
//
// 	The controller responds to the user input and performs interactions on the data model objects.
// 	The controller receives the input, optionally validates it and then passes the input to the model.
//
// To me when I played with the Separate Controller/UseCase layers in a production app, I felt the extra layer just was cumbersome.
// Therefore the "controller" concept here is a little different from traditional controllers but not different what most of the time a controller do in real world applications.
// Implementing business use cases, validation or you name it. From the simples case such as listing users, to the extend of executing heavy calculations on user behalf.
// The classical controller that provides interface between an external interface and a business logic implementation is something that must be implemented in the external interface integration layer itself.
// I removed this extra layer to make controller scope only to control the execution of a specific business use case based on generic inputs such as presenter and request.
// Also there is a return error which represent the unexpected and not recoverable errors that should notified back to the caller to teardown execution nicely.
//
//  The software in this layer contains application specific business rules.
//  It encapsulates and implements all of the use cases of the system.
//  These use cases orchestrate the flow of data to and from the entities,
//  and direct those entities to use their enterprise wide business rules to achieve the goals of the use case.
//
//  We do not expect changes in this layer to affect the entities.
//  We also do not expect this layer to be affected by changes to externalities such as the database,
//  the UI, or any of the common frameworks.
//  This layer is isolated from such concerns.
//
//  We do, however, expect that changes to the operation of the application will affect the use-cases and therefore the software in this layer.
//  If the details of a use-case change, then some code in this layer will certainly be affected.
//
//  Robert Martin
//
type Controller interface {
	Serve(Presenter, Request) error
}

// ControllerFunc helps convert anonimus lambda expressions into valid Frameless Controller object
type ControllerFunc func(Presenter, Request) error

// Serve func implements the Frameless Controller interface
func (lambda ControllerFunc) Serve(p Presenter, r Request) error {
	return lambda(p, r)
}

// Presenter is represent a communication layer presenting layer
//
// Scope:
// 	receive messages, and convert it into a serialized form
//
// You should not allow the users of the Presenter object to modify the state of the enwrapped communication channel, such as closing, or direct writing
type Presenter interface {
	//
	// RenderWithTemplate a content on a channel that the Presenter implements
	//	name helps to determine the what template should be used, but should not include channel specific names
	//	data is the content that should be used in the template
	// RenderWithTemplate(name string, data frameless.Content) error

	//
	// Render renders a simple message back to the enwrapped communication channel
	//	message is an interface type because the channel communication layer and content and the serialization is up to the Presenter to implement
	//
	// If the message is a complex type that has multiple fields,
	// an exported struct that represent the content must be declared at the controller level
	// and all the presenters must based on that input for they test
	Render(message interface{}) error
}

// PresenterFunc is a wrapper to convert standalone functions into a presenter
type PresenterFunc func(interface{}) error

// Render implements the Presenter Interface
func (lambda PresenterFunc) Render(message interface{}) error {
	return lambda(message)
}

// Request is framework independent way of interacting with a request that has been received on some kind of channel.
// from this, the controller should get all the data and options that should required for the business use case that use it.
type Request interface {
	Context() context.Context
	Data() Iterator
}

// Decoder is the interface for populating/replacing a public struct with values that retrived from an external resource
type Decoder interface {
	// Decode will populate an object with values and/or return error
	Decode(interface{}) error
}

// DecoderFunc enables to use anonimus functions to be a valid DecoderFunc
type DecoderFunc func(interface{}) error

// Decode proxy the call to the wrapped Decoder function
func (fn DecoderFunc) Decode(i interface{}) error {
	return fn(i)
}

// Iterator define a separate object that encapsulates accessing and traversing an aggregate object.
// Clients use an iterator to access and traverse an aggregate without knowing its representation (data structures).
// inspirated by https://golang.org/pkg/encoding/json/#Decoder
type Iterator interface {
	// this is required to make it able to cancel iterators where resource being used behind the schene
	// for all other case where the underling io is handled on higher level, it should simply return nil
	io.Closer
	// Next will ensure that Decode return the next item when it is executed
	Next() bool
	// Err return the cause if for some reason by default the More return false all the time
	Err() error
	// this is required to retrive the current value from the iterator
	Decoder
}

//
// Storing And Retrival

// ID is the serialized form of an identification entry that is provided by a given storage and than used for fetching entries from it.
// this is actually just a primitive string type, the type declaration is only meant to be for documentation purpose
type ID = string

// QueryUseCase is a storage specific component.
// It represent a use case for a specific Read of Update action.
// It should not implement or represent anything how the database query will be build,
// only contain necessary data to make it able implement in the storage.
// There must be no business logic located in a QueryUseCase type,
// this also includes complex types that holds logical information.
// It should only include as primitive fields as possible,
// and anything that require some kind of rule/logic to be done before it can be found in the storage should be done at controller level.
//
// By convention this should be declared where the Controllers are defined.
// If it's used explicitly by one Controller, the query use case definition should defined prior to that controller
// and listed close to each other in the go documentation.
//
type QueryUseCase = interface {
	// Test specify the given query Use Case behavior, and must be used for implementing it in the storage side.
	// For different context and scenarios testing#T.Run is advised to be used.
	// Test should create and teardown it's context in each execution (on each T.Run basis if T.Run used).
	//
	// The *testing.T is used, so the specification can use test specific methods like #Run and #Parallel.
	//
	Test(spec *testing.T, subject Storage)
}

// Storage define what is the most minimum that a storage should implement in order to be able
//
// The CRUD not enforced here because I found that most of the time, it is not even the case.
// For example, if there is no such query use case where you have to update something, why implement it ?
// I found that the most required functionality is Persisting an entity (it's a storage after all) and query and execute on it.
type Storage interface {
	io.Closer
	// Create able to create a given entity
	// By convention it also should set the id in the entity
	Create(Entity) error
	//
	// FindBy implements each application QueryUseCase.Test that used with the given storage.
	// This way for example the controller defines what is required in order to fulfil a business use case and storage implement that..
	//
	// This way it maybe feels boilerplated for really dynamic and complex searches, but for those,
	// I highly recommend to implement a separate structure that works with an iterator and do the complex filtering on the elements.
	// This way you can use easy to implement and minimalist query logic on the storage, and do complex things that is easy to test in a filter struct.
	//
	// By convention the QueryUseCase name should start with "[EntityName][FindLogicDescription]" so it is easy to distinguish it from other exported Structures,
	// 	example: UserByName, UsersByName, UserByEmail
	//
	// for simple use cases where returned iterator expected to only include 1 element I recommend using iterators.DecodeNext(iterator, &entity) for syntax sugar.
	// The use case common but I do not see benefit enforcing a First, FindOne or similar requirement from the storage.
	// For creating iterators for a single entity, you can use iterators.NewSingleElement.
	Find(QueryUseCase) Iterator
	//
	// Exec can execute a QueryUseCase that goal is to do modification in the storage in a described way by the QueryUseCase#Test method.
	//
	// By convention the QueryUseCase name should start with the Task to be achieved and than followed by type.
	// 	example: UpdateUserByName, DeleteUser, InvalidateUser
	Exec(QueryUseCase) error
}
