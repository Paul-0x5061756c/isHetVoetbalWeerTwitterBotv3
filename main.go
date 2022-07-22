package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)



func main(){
	lambda.Start(Handler)
}

func Handler(){
	u := <- getWeather();
	weatherStruct, err := UnmarshalWeatherStruct([]byte(u))

	if(err != nil){
		fmt.Println("something went wrong with fetching the weather")
	}

	payload := createTweet(weatherStruct)


	config := oauth1.NewConfig(os.Getenv("consumerKey"), os.Getenv("consumerKeySecret"))

	token := oauth1.NewToken(os.Getenv("accessToken"), os.Getenv("accessTokenSecret"))
	// http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Send a Tweet
	tweet, resp, err := client.Statuses.Update(payload, nil)

	fmt.Println(tweet, resp, err)


}
	
func getWeatherDescription(input string) string {
	m := make(map[string]string)
	m["Clouds"] = "bewolkt"
	m["Clear"] = "helder"
	m["Mist"] = "mistig"
	m["Snow"] = "aan het sneeuwen"
	m["Rain"] = "regenachtig"
	m["Drizzle"] = "miezerig"

	return m[input]
}

func canWePlayFootball(temp int, desc string) string {
	notSoNiceCondition := []string{"Rain", "Drizzle", "Snow"}

	idx := getIndexInArray(notSoNiceCondition, desc)

	if(temp > 20){
		if(idx == -1){
			return "Ja, maar het is wel erg warm"
		}else {
			return "Nee, het is warm, maar kans op neerslag"
		}
	}
	if(temp > 10){
		if(idx == -1){
			return "Ja, het is heerlijk"
		}else {
			return "Nee, temperatuur is lekker, maar kans op neerslag"
		}

	}else{
		if(idx == -1){
			return "Ja, maar wel wat fris"
		}else {
			return "Nee, het is houd & kans op neerslag"
		}
	}
}

func getIndexInArray(arr []string, needle string) int {
	for i, s := range arr {
    if(s == needle){
			return i
		}
	}
	return -1
}

func getWeather() <-chan string {
	url := os.Getenv("weatherUrl")
	r := make(chan string)

	go func() {
		defer close(r)
		res, _ := http.Get(url)
		text, _ := ioutil.ReadAll(res.Body)
		r <- string(text)
	}()
	return r
}

func createTweet(weather WeatherStruct) string{
	canWePlay := canWePlayFootball(int(weather.Main.Temp), weather.Weather[0].Main);
	desc := getWeatherDescription(weather.Weather[0].Main)
	loc, _ := time.LoadLocation("Europe/Amsterdam")
	now := time.Now().In(loc)

	return fmt.Sprintf("[%s] %s. De huidige temperatuur ligt rond de %.1f graden & het is voornamelijk %s. \n\n #voetbal #voetbalweer #football #weather #Levarne #ishetvoetbalweer`", now.Format("15:04") ,canWePlay, weather.Main.Temp, desc)

}

func UnmarshalWeatherStruct(data []byte) (WeatherStruct, error) {
	var r WeatherStruct
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *WeatherStruct) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type WeatherStruct struct {
	Coord      Coord     `json:"coord"`     
	Weather    []Weather `json:"weather"`   
	Base       string    `json:"base"`      
	Main       Main      `json:"main"`      
	Visibility int64     `json:"visibility"`
	Wind       Wind      `json:"wind"`      
	Clouds     Clouds    `json:"clouds"`    
	Dt         int64     `json:"dt"`        
	Sys        Sys       `json:"sys"`       
	Timezone   int64     `json:"timezone"`  
	ID         int64     `json:"id"`        
	Name       string    `json:"name"`      
	Cod        int64     `json:"cod"`       
}

type Clouds struct {
	All int64 `json:"all"`
}

type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Main struct {
	Temp      float64 `json:"temp"`      
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`  
	TempMax   float64 `json:"temp_max"`  
	Pressure  int64   `json:"pressure"`  
	Humidity  int64   `json:"humidity"`  
}

type Sys struct {
	Type    int64  `json:"type"`   
	ID      int64  `json:"id"`     
	Country string `json:"country"`
	Sunrise int64  `json:"sunrise"`
	Sunset  int64  `json:"sunset"` 
}

type Weather struct {
	ID          int64  `json:"id"`         
	Main        string `json:"main"`       
	Description string `json:"description"`
	Icon        string `json:"icon"`       
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int64   `json:"deg"`  
	Gust  float64 `json:"gust"` 
}

