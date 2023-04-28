pacakage main

import(
  "database/sql"
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  
  "github.com/google/uuid"
  "github.com/gorilla/mux"
  "_github.com/go-sql-driver/mysql"
)

type Todo struct{
  ID          int    `json:"id"`
  UUID        string `json:"uuid"`
  Description string `json:"description"`
  Status      bool   `json:"status"`
}

var db *sql.DB

func main() {
  dbHost :=os.Getenv("DB_HOSTNAME")
  dbUser :=os.Getenv("DB_USERNAME")
  dbPass :=os.Getenv("DB_PASSWORD")
  dbName :=os.Getenv("DB_NAME")
  
  var err error
  db, err=sql.Open("mysql",dbUser+":"+dbPass+"@tcp("+dbHost +":3306" +")/"+dbName)
  if err !=nil {
    log.Fatal(err)
  }
  defer db.Close()
  
  _,err=db.Exec("DROP TABLE IF EXISTS todos")
  if err!=nil {
    log.Fatal(err)
  }
  
  _,err = db.Exec("CREATE TABLE IF EXISTS todos (id INT AUTO_INCREMENT PRIMARY KEY,uuid VARCHAR(36) NOT NULL, description TEXT NOT NULL, status BOOLEAN NOT NULL DEFAULT FALSE);")
  if err!=nil {
    log.Fatal(err)
  }
  
  router :=mux.NewRouter()
  
  router.HandleFunc("/todos", createTodo).Methods("POST")
  router.HandleFunc("/todos/{id}", getTodoById).Methods("GET")
  router.HandleFunc("/todos", getAllTodos).Methods("GET")
  router.HandleFunc("/todos/{id}", updateTodo).Methods("PUT")
  router.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")
  log.Fatal(http.ListenAndServe(":8080",router))
  
 }
 func createTodo(w http.ResponseWriter, r *http.Request) {
    var todo Todo
    err :=json.NewDecoder(r.Body).Decode(&todo)
    if err!=nil {
      http.Error(w,err.Error(), http.StatusBadRequest)
      return
    }
    
    todo.UUID = uuid.New().String()
    todo.Status = true
    result,err :=db.Exec("INSERT INTO todos (uuid, description,status) VALUES (?,?,?)", todo.UUID, todo.Description, todo.Status)
    if err!=nil {
      http.Erroro(w,err.Error(),http.StatusInternalServerError)
      return
    }
    
    id,_ :=result.LastInsertId()
    todo.ID= int(id)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
  }
  
  func getTodoById(w http.ResponseWriter, r *http.Request) {
    params :=mux.Vars(r)
    id :=params["id"]
    
    var todo Todo
    err :=db.QueryRow("SELECT *FROM todos WHERE id=?", id).Scan(&todo.ID, &todo.UUID, &todo.Description, &todo.Status)
    if err==sql.ErrRows {
      http.NotFound(w,r)
      return
      }
      else if err!=nil {
        http.Error(w,err.Error(),http.StatusInternalServerError)
        return
      }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}

func getAllTodos(w http.ResponseWriter, r *http.Request) {
    rows,err :=db.QueryRow("SELECT *FROM todos")
    if err!=nil {
        http.Error(w,err.Error(),http.StatusInternalServerError)
        return
      }
      defer rows.Close()
      
      todos:=[]Todo{}
      for rows.Next(){
         var todo Todo
         err :=rows.Scan(&todo.ID, &todo.UUID, &todo.Description, &todo.Status)
         if err!=nil {
        http.Error(w,err.Error(),http.StatusInternalServerError)
        return
      }
      todos=append(todos, todo)
    }
    err=rows.Err()
    if err!=nil {
        http.Error(w,err.Error(),http.StatusInternalServerError)
        return
      }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
    params :=mux.Vars(r)
    id :=params["id"]
    
    var todo Todo
    err :=json.NewDecoder(r.Body).Decode(&todo)
    if err!=nil {
      http.Error(w,err.Error(), http.StatusBadRequest)
      return
    }
    
    todo.Status = true
    _,err =db.Exec("UPDATE todos SET description= ?, status=? WHERE id= ?", todo.Description, todo.Status, id)
     if err!=nil {
      http.Error(w,err.Error(), http.StatusBadRequest)
      return
    }
     w.Header().Set("Content-Type", "application/json")
     fmt.Fprintf(w,"Todo with ID %s was updated", id) 
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
    params :=mux.Vars(r)
    id :=params["id"]
    
    _,err :=db.Exec("DELETE FROM todos WHERE id= ?",id)
    if err!=nil {
      http.Error(w,err.Error(), http.StatusBadRequest)
      return
    }
     w.Header().Set("Content-Type", "application/json")
     fmt.Fprintf(w,"Todo with ID %s was Deleted", id) 
}
    

  
  
  
  
  
  
  
  
  
  
  
  
  
  
  
  
  
  
  
  
  
  
  
  
  
  
