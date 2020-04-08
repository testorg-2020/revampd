--------------------------------------------------------
--  DDL for Table UNIT_UNIVERSE
--------------------------------------------------------

CREATE TABLE "UNIT_UNIVERSE" 
   (	
  OP_YEAR INTEGER NOT NULL,
	UNIT_ID INTEGER NOT NULL, 
	EPA_REGION INTEGER NOT NULL, 
	STATE CHAR(2) NOT NULL, 
	FACILITY_NAME VARCHAR(120) NOT NULL, 
	ORIS_CODE INTEGER NOT NULL, 
	UNITID CHAR(6) NOT NULL, 
	STACK_IDS VARCHAR(200), 
	OP_STATUS CHAR(7) NOT NULL, 
	PROGRAM_CODE VARCHAR(250) NOT NULL, 
	UNIT_TYPE_DESCRIPTION VARCHAR(400) NOT NULL, 
	PRIMARY_FUEL_TYPE_DESC VARCHAR(200) NOT NULL, 
	PRIMARY_FUEL_GROUP VARCHAR(100) NOT NULL, 
	AN_COUNT_OP_TIME NUMERIC(12,0), 
	AN_GLOAD NUMERIC(12,0), 
	AN_SLOAD NUMERIC(12,0), 
	AN_HEAT_INPUT NUMERIC(12,1), 
	AN_CO2_MASS NUMERIC(12,1), 
	AN_SO2_MASS NUMERIC(12,1), 
	AN_NOX_MASS NUMERIC(12,1),
  CONSTRAINT year_unit_pk PRIMARY KEY(OP_YEAR,UNIT_ID)
   ) 
