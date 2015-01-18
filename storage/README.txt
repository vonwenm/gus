Storage drivers must follow these rules in order to work properly:

1. They must register themselves with the storage driver:
  func init() {
  	storage.Register(name_to_use_for_selection, &DriverStructure{})
  }

2. The driver (identified by the DriverStructure) must have one function: Open():

    func (t DriverStructure) Open(dsnConnect string, options string) (storage.Conn, error)

        The string 'options' is set in the configuration and passed to the open call. It can be ignored or
        provide extra information to use within the driver.

    The 'storage.conn' is an interface that must implement the following calls:

    	UserUpdate(user *record.User) error
    	UserInsert(user *record.User) error

    	UserFetch(key, value string) (*record.User, error)


    It may optionally implement:
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

3. All functions from all classes must return errors of the type defined by ecode.ErrorCoder. If you want
    to return additional information, for example status or field information, you should create another interface
    that implements the same as ErrorCoder but with additional fields.

4. The record.User level manipulates the in-memory image of a user. It needs to perform any operations that
    will alter or set information for each service call.

5. The store-level driver has very few operations: Insert, Update and Fetch. Records are not deleted by GUS but can be
    marked for deletion instead by the record/User level. The driver is only responsible for mapping the record-level
    fields back and forth with the database level fields.

6. The driver-level interface has 'aliases' for Fetches (e.g. UserFetchByGuid) that make calling the lower routines
    a little bit easier.

