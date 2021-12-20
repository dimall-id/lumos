
# Lumos

Lumos is a collection of library, wrapper and standard library to build a service in Dimall Backend Project.

Main Feature :
 1. Routing
 2. Data Query

## Add Routing

    import web "github.com/dimall-id/lumos/http"
    import "net/http"
    
    func Routes () {
	    err := web.AddRoute(http.Route{  
			  Name: "List of Product",  
			  HttpMethod: "GET",  
			  Url: "/products",  
			  Roles: []string{  
			      "USER",  
			  },  
			  Func: ListOfProductHandler,
			}
		)  
		if err != nil {  
		   log.Fatal(err)  
		}
    }
    
    func ListOfProductHandler (r *http.Request) (interface{},web.HttpError) {
	    return nil,web.HttpError{}
    }


## Start Web Server

    import web "github.com/dimall-id/lumos/http"
    
    func main() {
	    Routes()
	    log.Fatal(web.StartHttpServer(":8080"))
    }

## Query Data

 1. Date Query
	 format : field=[gt|gte|eq|neq:mm-dd-yyyy,lt|lte:mm-dd-yyyy]
	 example :
	 - date=[gt:01-20-2020,lt:01-22-2020]
	 - date=[eq:01-20-2020]
 2. Numeric Query
	 format : field=[gt|gte|eq|neq:0-9,lt|lte:0-9]
	 example :
	 - price=[gt:1000,lte:10000]
	 - price=[neq:100000]
 3. List Query
	 format : field=[in|nin:value,value,...]
	 example :
	 - id=[in:1,2,3,4]
	 - id=[nin:1,2,3,4]
 4. Order Query
	 format : order=[field:asc|desc,field:asc|desc,...]
	 example :
	 - order=[name:asc,price:desc]
 5. Paging Query
	 format : paging=[page:0-9,per_page:0-9]
	 example :
	 - paging=[page:1,per_page:10]
 6. With Query
	 format : with=[with:relation,relation,...]
	 example :
	 - with=[with:ProductImage,ProductCategories]
 7. Select Query
	 format : select=[select:field,field,field,...]
	 example :
	 - select=[select:id,name,price]
 8. String Query
	 format : field=[eq|neq|like|ilike:value]
	 example
	- name=[eq:Product A]
	- name=[neq:Product B]
	- name=[like:Produc]
	- name=[ilike:Product]


## Setup Query at Handler

    import "gorm.io/gorm"
    import "github.com/dimall-id/lumos/data"
    func ListOfProductHandler (r *http.Request) (interface{},web.HttpError) {
	    tx := db.Open(...)
	    tx.Model(Product)
	    
	    Q = data.New(tx)
	    
	    var data []Product
	    response := Q.BuildResponse(r, &data)
	    
	    return response, http.HttpError{}
    } 
