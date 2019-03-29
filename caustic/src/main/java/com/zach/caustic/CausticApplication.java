package main.java.com.zach.caustic;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.http.client.support.BasicAuthorizationInterceptor;

@SpringBootApplication
public class CausticApplication {

    public static void main(String[] args) {
        SpringApplication.run(CausticApplication.class, args);
    }

}
