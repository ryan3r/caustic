package com.ryan3r.caustic.controller;
import java.io.IOException;

import com.ryan3r.caustic.repository.submitsRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

@RestController
public class submitsController {
	@Autowired
	submitsRepository s;
	
	@RequestMapping(method=RequestMethod.POST, value="/submit", consumes=MediaType.MULTIPART_FORM_DATA_VALUE)
	public void uploadMapServer(@RequestParam("upload") MultipartFile file) throws IOException {
		
	}
	
	
}
