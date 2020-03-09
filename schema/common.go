package schema

import (
	"time"
)

type Data struct {
	ID         string      
	Channel    chan string 
	Active     *bool       
	Status     *string     
	StepTime   int         
	Start      int         
	Count      *int        
	CreatedAt  time.Time   
	ModifiedAt time.Time   
}

type CacheData struct {
	ID         string      
	Active     *bool       
	Status     *string     
	StepTime   int         
	Start      int         
	Count      *int        
	CreatedAt  time.Time   
	ModifiedAt time.Time   
}

type CheckData struct {
	ID        string    
	StepTime  int       
	Count     int       
	Status    string    
	CreatedAt time.Time 
}

type PauseData struct {
	ID        string    
	PauseTime time.Time 
}

type Response struct {
	Data    interface{} 
	Success bool        
}

type HTMLRender struct {
	Render []CheckData
}
