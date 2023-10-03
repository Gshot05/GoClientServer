CREATE TABLE Type_Display(
	ID_Type_Display serial not null constraint PK_Type_Display primary key,
	Name_Diagonal decimal(4, 1) not null,
	Name_Resolution varchar(70) not null,
	Type_Type varchar(70) not null,
	Type_Gsync boolean not null
);

CREATE TABLE Type_Monitor(
	ID_Type_Monitor serial not null constraint PK_Type_Monitor primary key,
	Name_Voltage decimal(4, 1) not null,
	Name_Gsync_Prem boolean not null,
	Name_Curved boolean not null,
	Type_Display_ID int not null references Type_Display(ID_Type_Display)
);
select * from Type_Display;
select * from Type_Monitor;

drop TABLE Type_Monitor;
drop TABLE Type_Display;


CREATE TABLE Type_Users(
	ID_Type_Users serial not null constraint PK_Type_Users primary key,
	Name_Username text not null constraint UQ_Type_Users Unique,
	Name_Password text not null,
	Name_email text not null,
	Name_Is_Admin boolean not null
);
INSERT INTO Type_Users (ID_Type_Users, 
						Name_Username, 
						Name_Password, 
						Name_email, 
						Name_Is_Admin)
values (1, 
		'Admin', 
		'42b1c2497694ae761ba3f36ea19b735e1f6fb77949ad676b01a2758dfcf05285', 
	    'email@mail.com', 
	    true);
select * from Type_Users;
