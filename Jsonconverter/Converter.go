package Jsonconverter

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)
//Go-datatypes
var (
	s_struct      = " struct{\n"
	s_array       = " []"
	s_string      = " string\n"
	s_int         = " int\n"
	s_bool        = " bool\n"
	s_interface   = " interface{}\n"
	s_close_curly = "\n}\n"
)
//Checks for different conditions
var arr_alert bool
var arr_struct_alert bool
var garbage bool
var arr_var_check bool
var arr_check bool
var arr_may_end bool
var close_arr bool
var datastore []string

var mystruct string

func ValidateJson(c *gin.Context) {
	databyte, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	content := string(databyte)
	content = strings.TrimSpace(content)
	_ = ioutil.WriteFile("json_example_write.txt", []byte(content), 0644)
	res1 := strings.HasPrefix(content, `{`)
	res2 := strings.HasSuffix(content, `}`)
	if !res1 || !res2 {
		fmt.Println("Error Improper Json")
		return
	}
	opencurlycount := strings.Count(content, `{`)
	closecurlycount := strings.Count(content, `}`)
	totalcur := opencurlycount - closecurlycount
	openarrbrac := strings.Count(content, `[`)
	closedarrbrac := strings.Count(content, `]`)
	totalarrbrac := openarrbrac - closedarrbrac
	if totalarrbrac != 0 || totalcur != 0 {
		fmt.Println("Error Improper Json")
		return
	}
	packagename := c.Query("PackageName")
	StructName := c.Query("StructName")

	mystruct = "package " + packagename + "\n\ntype " + StructName + " struct{\n"
	Converter()
	c.JSON(200, "OK")

}
func ReadFile() []string {
	file, err := os.Open("json_example_write.txt")
	if err != nil {
		log.Fatalf("failed to open")

	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	file.Close()
	return text
}
func Converter() {
	text := ReadFile()

	for i := 1; i < len(text); i++ {
		fmt.Println("----->first occurence = ", arr_check, text[i])
		if strings.Contains(text[i], ":") && !arr_alert {
			fmt.Println("\nHHHHHHH = =", text[i])
			split := strings.Split(text[i], ":")
			key := split[0]
			key = IsLetter(key)
			Value := split[1]
			Value = checkdatatype(Value)
			mystruct += key + Value
		}else if (arr_check) && strings.Contains(text[i], "]") {
			fmt.Println("dddddddd == ", datastore)
				arr_check = false
				var data string
				for i := 0; i < len(datastore); i++ {
					if datastore[i] == datastore[0] {
						data = datastore[0]
						fmt.Println("aaaaa = ",data) 
					} else {
						data = s_interface
						fmt.Println("aaaaa = ",data) 
					}
				}
				mystruct += data
				datastore = nil
			arr_var_check, arr_alert, arr_var_check, garbage = false, false, false, false
		} else if !arr_alert {
			if strings.Contains(text[i], `}`) {
				mystruct += s_close_curly
				fmt.Println("\niiiiiiii = =", text[i])
			} //for array inside struct
		} else if arr_alert && !arr_struct_alert && !arr_var_check && !garbage {
			//array containing struct
			if strings.Contains(text[i], `{`) {
				arr_struct_alert = true
				mystruct += s_struct
				fmt.Println("\njjjjjjjjjj = =", text[i])
			} else {
				fmt.Println("qqqqqqq == ", text[i])
				arr_var_check = true
				arr_check = true
				data := checkdatatype(text[i])
				datastore = append(datastore, data)
			}
			fmt.Println(arr_struct_alert, arr_var_check)
		} else if !garbage && !arr_var_check && arr_struct_alert {
			fmt.Println("=!!!!!!!", text[i])
			if strings.Contains(text[i], `[`) {
				arr_var_check = true
			} 
			if strings.Contains(text[i], ":") {
				fmt.Println("inside = ", text[i])
				split := strings.Split(text[i], ":")
				key := split[0]
				key = IsLetter(key)
				Value := split[1]
				Value = checkdatatype(Value)
				mystruct += key + Value
			} else if strings.Contains(text[i], `}`) {
				fmt.Println("struct_array", text[i])
				// mystruct += s_close_curly
				garbage = true
				arr_struct_alert = false
			}
		} else if arr_var_check && arr_alert && arr_struct_alert && !strings.Contains(text[i], `]`) {
			val := strings.TrimSpace(text[i])
			data := checkdatatype(val)
			datastore = append(datastore, data)
			fmt.Println("know data", val, data, text[i])
		} else if strings.Contains(text[i], `]`) && !garbage {
			// arr_alert = false
			fmt.Println("got dtata = ", text[i])
			var data string
			for i := 0; i < len(datastore); i++ {
				if datastore[i] == datastore[0] {
					data = datastore[0]
				} else {
					data = s_interface
				}
			}
			mystruct += data
			datastore = nil
			// fmt.Println(datastore)
			// fmt.Println("close end",text[i])
			arr_var_check = false
		} else if (garbage) && strings.Contains(text[i], `}`) && !strings.Contains(text[i], `,`) {
			fmt.Println("garbage = ", text[i])
			arr_may_end = true
			// mystruct += text[i]+"\n"
			// arr_var_check, arr_alert, arr_var_check, garbage = false, false, false, false
		} else if arr_var_check && arr_check {
			fmt.Println("going", text[i])
			val := strings.TrimSpace(text[i])
			data := checkdatatype(val)
			datastore = append(datastore, data)
			fmt.Println("getmydata", datastore)
			
		}else if (arr_may_end){
			if strings.Contains(text[i],`]`){
				close_arr = true
				mystruct += s_close_curly+"\n"
			arr_var_check, arr_alert, arr_var_check, garbage,arr_check,arr_may_end = false,false,false, false, false, false
			}else{
				fmt.Println("OOPs! it's not an end")
			}
		}
	}
	fmt.Println(mystruct)
	WriteFile(mystruct)

}
func IsLetter(s string) string {
	str := ""
	s = strings.TrimSpace(s)

	for _, r := range s {
		if unicode.IsLetter(r) {
			str += string(r)
		}

	}
	if str != "" {
		str = strings.Replace(str, string(str[0]), strings.ToUpper(string(str[0])), 1)
	}
	return str
}
func checkdatatype(s string) string {
	str := ""
	s = strings.TrimSpace(s)
	datacheck := string(s[0])
	if strings.Contains(datacheck, `"`) {
		str = s_string
	} else if regexp.MustCompile(`\d`).MatchString(s) {
		str = s_int
	} else if strings.Contains(datacheck, "t") || strings.Contains(datacheck, "f") {
		str = s_bool
	} else if strings.Contains(datacheck, `[`) {
		str = s_array
		arr_alert = true
	} else if strings.Contains(datacheck, `{`) {
		str += s_struct
	}
	return str
}
func WriteFile(string){
	_ = os.WriteFile("Go_example_struct.go", []byte(mystruct), 0644)
}