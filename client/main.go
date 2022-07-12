package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gofiber/fiber/v2"
)


func main () {
	app := fiber.New()

	app.Get("/api", api)
	app.Get("/api2", api2)

	app.Listen(":8001")
}

func init () {

	hystrix.ConfigureCommand("api", hystrix.CommandConfig{
		Timeout: 				500,
		RequestVolumeThreshold: 4,  // Test : Config request to open circuit breaker
		ErrorPercentThreshold:  50, // Default 50 % 
		SleepWindow: 			15000, // Open circuit breaker 15 s
	})

	hystrix.ConfigureCommand("api2", hystrix.CommandConfig{
		Timeout: 				500,
		RequestVolumeThreshold: 4,  // Test : Config request to open circuit breaker
		ErrorPercentThreshold:  50, // Default 50 % 
		SleepWindow: 			15000, // Open circuit breaker 15 s
	})

	// ทำ dashboard
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(":8002", hystrixStreamHandler) // ใช้ goroutine ในการ start

	//docker pull mlabouardy/hystrix-dashboard
	//สร้าง file docker-compose.yml
	//docker compose up -d
	// http://localhost:9002/hystrix
}

func api(c *fiber.Ctx) error {

	output := make(chan string, 1)
	hystrix.Go("api", func () error {
		res, err := http.Get("http://localhost:8000/api")
		if err != nil {
			return err
		}
		defer res.Body.Close()
	
		data, err := io.ReadAll(res.Body) // call api ปกติ
		if err != nil {
			return err
		}
	
		msg := string(data)
		fmt.Println(msg)

		output <- msg

		return nil
	}, func (err error) error {
		fmt.Println(err)
		return nil
	})

	out := <- output

	return c.SendString(out)
}

func api2(c *fiber.Ctx) error {

	output := make(chan string, 1)
	hystrix.Go("api2", func () error {
		res, err := http.Get("http://localhost:8000/api")
		if err != nil {
			return err
		}
		defer res.Body.Close()
	
		data, err := io.ReadAll(res.Body) // call api ปกติ
		if err != nil {
			return err
		}
	
		msg := string(data)
		fmt.Println(msg)

		output <- msg

		return nil
	}, func (err error) error {
		fmt.Println(err)
		return nil
	})

	out := <- output

	return c.SendString(out)
}