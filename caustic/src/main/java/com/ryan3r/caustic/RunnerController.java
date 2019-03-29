package com.ryan3r.caustic;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class RunnerController {
	@RequestMapping(method = RequestMethod.GET, path="/ok")
	public String slash() {
		return "OK";
	}
}
