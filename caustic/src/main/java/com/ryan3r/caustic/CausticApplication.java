package com.ryan3r.caustic;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class CausticApplication {
	public static long contestStartTime = System.currentTimeMillis() / 60000; // in minutes

	public static void main(String[] args) {
		SpringApplication.run(CausticApplication.class, args);
	}

}
