package memory


func  (t *StorageMem) CreateStore() error {
	if err:= t.createUser();err != nil {
	return err
	}
	return t.createUserData()
}


func (t *StorageMem)  createUserData( ) error {
	sql := `
        CREATE TABLE IF NOT EXISTS UserData(
        	Guid text ,
        	Name text ,
        	Type text ,
        	Encode text,
        	Data blob ,
        	Primary KEY ( Guid, Name, Type ) );
        create index idx_name on UserData( Name );`
	_,err :=t.db.Exec(sql)
	return err
}

func  (t *StorageMem) createUser() error {
	sql := `
	CREATE TABLE IF NOT EXISTS User (
		Guid text primary key,
		FullName text,
		Email    text,


		Domain    text,
		LoginName text,
		Password  text,
		Token     text,

		Salt text,

		IsActive   integer,
		IsLoggedIn integer,

		LoginAt      text,
		LogoutAt     text,
		LastAuthAt   text,
		LastFailedAt text,
		FailCount    integer ,

		MaxSessionAt text,
		TimeoutAt  text,

		CreatedAt text,
		UpdatedAt text,
		DeletedAt text);
	CREATE unique index uemail on User(Email);
	CREATE unique index ulname on User(LoginName);
	CREATE index ufullname on User(FullName);
	CREATE index utoken on User(Token);
	`
	_, err := t.db.Exec(sql)
	return err
}
