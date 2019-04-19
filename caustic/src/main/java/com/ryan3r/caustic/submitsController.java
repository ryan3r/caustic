package com.ryan3r.caustic;
import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

@RestController
public class submitsController {
	@Autowired
	submitsRepository s;
	
	@RequestMapping(method=RequestMethod.POST, value="/maps/{id}/upload/{user}", consumes=MediaType.MULTIPART_FORM_DATA_VALUE)
	public void uploadMapServer(@RequestParam("file") MultipartFile file, @PathVariable String id, @PathVariable String user) throws IOException {
		
	}
	
	
}
