package com.ryan3r.caustic.controller;

import com.ryan3r.caustic.repository.submitsRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;

@RestController
public class submitsController {
	@Autowired
	submitsRepository s;
	
	@RequestMapping(method=RequestMethod.POST, value="/submit", consumes=MediaType.MULTIPART_FORM_DATA_VALUE)
	public void uploadMapServer(@RequestParam("upload") MultipartFile file) throws IOException {
		String idp = new String("IDPLACEHOLDER");
		File f = new File("IDPLACEHOLDER","/mnt/submissions");
		if(!f.exists())
		{
			f.mkdir();
		}
		File fNew = new File(file.getName(), "/mnt/submissions/" + idp);
		fNew.createNewFile();
		FileOutputStream fout = new FileOutputStream(fNew);
		fout.write(file.getBytes());
		fout.close();
	}
	
	
}
