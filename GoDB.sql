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