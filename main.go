package main

import (
)

func main()  {

	//ip := flag.String("ip", "not-a-valid-ip", "self ip")
	//ips_str := flag.String("ips", "not-a-valid-ip", "ip1,ip2,ip3....")
	//flag.Parse()
	//ips := strings.Split(*ips_str, ",")
	//
	//service := Service{}
	//service.Start(*ip, ips)

	// cluster
	test1_ip := "127.0.0.1:11111"
	test1_ips := [] string{"127.0.0.1:22222", "127.0.0.1:33333"}
	service1 := Service{}
	go func() {
		service1.Start(test1_ip, test1_ips)
	}()

	test2_ip := "127.0.0.1:22222"
	test2_ips := [] string{"127.0.0.1:11111", "127.0.0.1:33333"}
	service2 := Service{}
	go func() {
		service2.Start(test2_ip, test2_ips)
	}()

	test3_ip := "127.0.0.1:33333"
	test3_ips := [] string{"127.0.0.1:11111", "127.0.0.1:22222"}
	service3 := Service{}
	go func() {
		service3.Start(test3_ip, test3_ips)
	}()

	for {
		select {
		}
	}
}