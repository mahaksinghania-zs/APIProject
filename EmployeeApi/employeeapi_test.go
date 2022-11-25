package main

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	uuid2 "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mock sqlmock.Sqlmock
var err error

func Test_GetEmployeeDetails(t *testing.T) {

	Db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer Db.Close()
	rows := sqlmock.NewRows([]string{"deptid", "deptName", "id", "name", "Phone_no"}).AddRow("6c69ce1c-6be2-11ed-9e01-64bc589457a0", "engineering", "fbd13799-6bd3-11ed-9e01-64bc589457a0", "MAHAK ", "245262728").AddRow("6c69ce1c-6be2-11ed-9e01-64bc589457a0", "technical", "fbd13799-6bd3-11ed-9e01-64bc589457a0", "riya", "2738398")

	mock.ExpectQuery("SELECT department.Id, department.Name ,employee.Id, employee.Name,employee.Phone FROM employee INNER JOIN department ON employee.DepartmentId=department.Id;").WillReturnRows(rows)
	testcases := []struct {
		expectedoutput string
		descr          string
		statusCode     int
	}{
		{

			expectedoutput: `[{"deptDetails":{"deptid":"6c69ce1c-6be2-11ed-9e01-64bc589457a0","deptName":"engineering"},"id":"fbd13799-6bd3-11ed-9e01-64bc589457a0","name":"MAHAK ","phone_no":"245262728"}]`,

			descr:      "getting the details of the employee",
			statusCode: 200,
		},
		{

			expectedoutput: `[{"deptDetails":{"deptid":"6c69ce1c-6be2-11ed-9e01-64bc589457a0","deptName":"technical"},"id":"fbd13799-6bd3-11ed-9e01-64bc589457a0","name":"MAHAK ","phone_no":"245262728"}]`,

			descr:      "getting the details of the employee",
			statusCode: 400,
		},
	}

	for _, v := range testcases {
		// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("GET", "/employees", nil)
		if err != nil {
			t.Errorf(err.Error())
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		response := httptest.NewRecorder() //response.body

		GetEmployeeDetails(response, req)
		//var val []Employee                              //just a decleration
		//_ = json.Unmarshal(response.Body.Bytes(), &val) //copying response into val//deseralization //converted from json to  go structure //deserialization

		assert.Equal(t, v.statusCode, response.Code) //it asserts whether two objects are equal or not.
		assert.Equal(t, v.expectedoutput, response.Body.String())

	}
}

func Test_GetEmployeeDetailsById(t *testing.T) {
	uuid := uuid2.NewString()
	//uuid := "6c69ce1c-6be2-11ed-9e01-64bc589457a0"
	Db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer Db.Close()
	rows := sqlmock.NewRows([]string{"deptid", "deptName", "id", "name", "Phone_no"}).AddRow("6c69ce1c-6be2-11ed-9e01-64bc589457a0", "engineering", uuid, "MAHAK ", "245262728")

	mock.ExpectQuery("SELECT department.Id, department.Name ,employee.Id, employee.Name,employee.Phone FROM employee INNER JOIN department ON employee.DepartmentId=department.Id where employee.Id=?").WithArgs(uuid).WillReturnRows(rows)
	testcases := []struct {
		expectedoutput string
		descr          string
	}{
		{

			expectedoutput: `{"deptDetails":{"deptid":"6c69ce1c-6be2-11ed-9e01-64bc589457a0","deptName":"engineering"},"id":"` + uuid + `","name":"MAHAK ","phone_no":"245262728"}`,
			//	expectedoutput: `[{"deptDetails":{"deptid":6c69ce1c-6be2-11ed-9e01-64bc589457a0",
			//	"deptName":engineering",
			//},
			//	"id":fbd13799-6bd3-11ed-9e01-64bc589457a0",
			//	"name":MAHAK",
			//	"Phone_no":245262728",
			//]`
			descr: "getting the details of the employee",
		},
	}

	for _, v := range testcases {

		req, err := http.NewRequest("GET", "/employees?id="+uuid, nil)
		if err != nil {
			t.Errorf(err.Error())
		}
		//id:=req.URL.Query().Set("id",uuid)

		response := httptest.NewRecorder() //response.body

		GetEmployeeDetailsById(response, req)
		//var val []Employee                              //just a decleration
		//_ = json.Unmarshal(response.Body.Bytes(), &val) //copying response into val//deseralization //converted from json to  go structure //deserialization

		//assert.Equal(t, v.statusCode, response.Code) //it asserts whether two objects are equal or not.
		assert.Equal(t, v.expectedoutput, response.Body.String())

	}
}

func Test_GetDepartmentDetailsById(t *testing.T) {
	uuid := uuid2.NewString()
	//uuid := "6c69ce1c-6be2-11ed-9e01-64bc589457a0"
	Db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer Db.Close()
	rows := sqlmock.NewRows([]string{"deptid", "deptName"}).AddRow(uuid, "engineering")

	mock.ExpectQuery("SELECT * FROM  department WHERE department.Id =?").WithArgs(uuid).WillReturnRows(rows)
	testcases := []struct {
		expectedoutput string
	}{
		{

			`{"deptid":"` + uuid + `","deptName":"engineering"}`,
		},
	}

	for _, v := range testcases {

		req, err := http.NewRequest("GET", "/employees?id="+uuid, nil)
		if err != nil {
			t.Errorf(err.Error())
		}
		//id:=req.URL.Query().Set("id",uuid)

		response := httptest.NewRecorder() //response.body

		GetDepartmentDetailsById(response, req)
		//var val []Employee                              //just a decleration
		//_ = json.Unmarshal(response.Body.Bytes(), &val) //copying response into val//deseralization //converted from json to  go structure //deserialization

		//assert.Equal(t, v.statusCode, response.Code) //it asserts whether two objects are equal or not.
		assert.Equal(t, v.expectedoutput, response.Body.String())

	}
}

func Test_GetDepartmentDetails(t *testing.T) {

	Db, mock, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer Db.Close()
	rows := sqlmock.NewRows([]string{"deptid", "deptName"}).AddRow("6c69ce1c-6be2-11ed-9e01-64bc589457a0", "engineering")

	mock.ExpectQuery("SELECT * from department;").WillReturnRows(rows)
	testcases := []struct {
		expectedoutput string
		descr          string
	}{
		{

			expectedoutput: `[{"deptid":"6c69ce1c-6be2-11ed-9e01-64bc589457a0","deptName":"engineering"}]`,
			//	expectedoutput: `[{"deptDetails":{"deptid":6c69ce1c-6be2-11ed-9e01-64bc589457a0",
			//	"deptName":engineering",
			//},
			//	"id":fbd13799-6bd3-11ed-9e01-64bc589457a0",
			//	"name":MAHAK",
			//	"Phone_no":245262728",
			//]`
			descr: "getting the details of the employee",
		},
	}

	for _, v := range testcases {
		// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("GET", "/depts", nil)
		if err != nil {
			t.Errorf(err.Error())
		}

		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		response := httptest.NewRecorder() //response.body

		GetDepartmentDetails(response, req)
		//var val []Employee                              //just a decleration
		//_ = json.Unmarshal(response.Body.Bytes(), &val) //copying response into val//deseralization //converted from json to  go structure //deserialization

		//assert.Equal(t, v.statusCode, response.Code) //it asserts whether two objects are equal or not.
		assert.Equal(t, v.expectedoutput, response.Body.String())

	}
}

func Test_CreateEmployee(t *testing.T) {

	Db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer Db.Close()
	//rows := sqlmock.NewRows([]string{"deptid", "deptName"}).AddRow("6c69ce1c-6be2-11ed-9e01-64bc589457a0", "engineering")

	mock.ExpectExec("insert into employee").WithArgs(sqlmock.AnyArg(), "fbd13799-6bd3-11ed-9e01-64bc589457a0", "MAHAK ", "245262728").
		WillReturnResult(sqlmock.NewResult(1, 1))

	//mock.ExpectQuery("insert into employee (Id, NAME,DepartmentID,PHONE) values (,?,?,?)\", emp.Name, emp.DeptDetails.DeptId, emp.PhoneNo").WillReturnRows(rows)

	testcases := []struct {
		description    string
		input          Employee
		expectedOutput Employee //can get more than one value at a time so using slice
		statusCode     int
	}{
		{"All entries are present",
			Employee{
				Department{
					"6c69ce1c-6be2-11ed-9e01-64bc589457a0",
					"engineering",
				},
				"fbd13799-6bd3-11ed-9e01-64bc589457a0", "Mahak",
				"245262728",
			},
			Employee{
				Department{
					"6c69ce1c-6be2-11ed-9e01-64bc589457a0",
					"engineering",
				},
				"fbd13799-6bd3-11ed-9e01-64bc589457a0", "Mahak",
				"245262728",
			},
			201,
		},
	}

	for _, tc := range testcases {
		val, _ := json.Marshal(tc.input) //go to json
		req, err := http.NewRequest("POST", "/employeee", bytes.NewReader(val))
		if err != nil {
			t.Errorf(err.Error())
		}
		//response recorder
		response := httptest.NewRecorder()
		CreateEmployee(response, req)
		var actRes Employee
		_ = json.Unmarshal(response.Body.Bytes(), &actRes) //json to go
		assert.Equal(t, tc.statusCode, response.Code)
		t.Error(actRes)
		//assert.Equal(t, tc.expectedOutput, actRes)

	}

}

func TestCreateDepartment(t *testing.T) {
	Db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer Db.Close()
	//rows := sqlmock.NewRows([]string{"deptid", "deptName"}).AddRow("6c69ce1c-6be2-11ed-9e01-64bc589457a0", "engineering")

	mock.ExpectBegin()
	mock.ExpectExec("insert into employee").WithArgs(sqlmock.AnyArg(), "engineering", "fbd13799-6bd3-11ed-9e01-64bc589457a0", "MAHAK ", "245262728").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	//mock.ExpectQuery("insert into employee (Id, NAME,DepartmentID,PHONE) values (,?,?,?)\", emp.Name, emp.DeptDetails.DeptId, emp.PhoneNo").WillReturnRows(rows)

	testcases := []struct {
		description    string
		input          Department
		expectedOutput Department //can get more than one value at a time so using slice
		statusCode     int
	}{
		{"All entries are present",
			Department{
				"6c69ce1c-6be2-11ed-9e01-64bc589457a0",
				"engineering",
			},

			Department{
				"6c69ce1c-6be2-11ed-9e01-64bc589457a0",
				"engineering",
			},
			200,
		},
	}

	for _, tc := range testcases {
		val, _ := json.Marshal(tc.input) //go to json
		req, err := http.NewRequest(http.MethodPost, "/department", bytes.NewBuffer(val))
		if err != nil {
			t.Errorf(err.Error())
		}
		//response recorder
		resRec := httptest.NewRecorder()
		CreateDepartment(resRec, req)
		var actRes Department
		_ = json.Unmarshal(resRec.Body.Bytes(), &actRes) //json to go
		assert.Equal(t, tc.statusCode, resRec.Code)
		//assert.Equal(t, tc.expectedOutput, actRes)

	}
}
