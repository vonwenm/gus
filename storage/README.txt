Storage drivers must follow these rules in order to work properly:

1. They must register themselves with the storage driver:
  func init() {
  	storage.Register(STORAGE_IDENTITY, &DriverStructure{})
  }

2. The driver (identified by the DriverStructure) must have one function: Open():
    func (t DriverStructure) Open(dsnConnect string) (storage.Conn, error)

    The 'storage.conn' is an interface that must implement the following calls:
        RegisterUser(user *record.User) error

    	UserLogin(user *record.User) error
    	UserAuthenticated(user *record.User) error
    	UserLogout(user *record.User) error

    	UserUpdate(user *record.User) error

    	FetchUserByGuid(guid string) (*record.User, error)
    	FetchUserByToken(token string) (*record.User, error)
    	FetchUserByEmail(email string) (*record.User, error)
    	FetchUserByLogin(loginName string) (*record.User, error)

    It may also implement:
        CreateStore() error
            If implemented, it should create any directory/resource/tables required
            for the driver to save information. If possible, it should be non-destructive, i.e.,
            it should not delete or destroy anything in the system.
        Close() error
            Close will allow you to close any connection to the database/file that has been created.
            Close should not return an error if it is called multiple times and there is nothing to
            close. (Meaning, keep the state in the driver and don't force service routines to track
            if it needs to be closed or not.)
        Reset()
            Reset should reset any errors and cleanup any data, for example if you have implemented any
            cursor reads. Calling Reset will never return an error.
        Ping() error
            Ping can be used to see if the data store connection is alive. Ping should only return
            an error if there is a data store connection error.
       	Release() error
       	    Release should release any locks/resources that have been set. This is useful for file locks
       	    or database record locks. Note that, if implemented, the function MUST be able to be called
       	    any number of times even if no lock/resources need to be released. The state of resources to
       	    be released must be kept by the driver, not by the caller.

3. All functions from all classes must return errors of the type defined by ecode/ErrorCoder. If you want
    to return additional information, for example status or field information, you should create another interface
    that implements the same as ErrorCoder but with additional fields.

4. The record/User level manipulates the in-memory image of a user. It needs to perform any operations that
    will alter or set information.
    The driver level should save and fetch records and update fields in the
    user record that reflect its operation. For example, if a login is successful, it will set the IsLoggedIn
    to true and the LoginAt time fields. This simplifies the logic in the service layer. If the save
    isn't successful, don't alter fields in the user record. (See the drivers included)

5. The reason for the different update/create calls is they each have different criteria to be met for their
    operation. This simplifies the upper level logic. The Fetch routines are simple lookups and
    each of the operations should be implemented. Each of the fetch routines are selecting by unique keys.
