package com.ryan3r.caustic.controller;

import com.ryan3r.caustic.model.Submission;
import com.ryan3r.caustic.repository.SubmissionRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.nio.file.Paths;

@RestController
public class submitsController {
	@Autowired
	SubmissionRepository repo;
	
	@PostMapping("/submit")
	public String uploadMapServer(@RequestParam("problemId") String problemId, @RequestParam("upload") MultipartFile file) throws IOException {
		Submission submission = new Submission(problemId, file.getOriginalFilename());
		submission = repo.save(submission);
		String id = "" + submission.getSubmissionId();

		File f = Paths.get("/mnt/submissions/", id).toFile();
		if(!f.exists())
		{
			f.mkdirs();
		}
		
		File fNew = submission.getFile();
		fNew.createNewFile();
		FileOutputStream fout = new FileOutputStream(fNew);
		fout.write(file.getBytes());
		fout.close();

		return "Ok";
	}
}
