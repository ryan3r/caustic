package com.ryan3r.caustic.controller;

import com.ryan3r.caustic.repository.accountsRepository;
import com.ryan3r.caustic.model.accounts;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.servlet.view.RedirectView;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.nio.file.Paths;

@RestController
public class accountsController 
{
	@Autowired
	accountsRepository a;

	@RequestMapping(method=RequestMethod.POST, value="/accounts")
	public boolean addAccount(@RequestBody accounts account)
	{
		if(a.findUser(account.getUsername()) != null)
			return false;
		a.save(account);

		return true;
	}
	@RequestMapping(method=RequestMethod.POST, value="/accountsLogin")
	public boolean loginAccount(@RequestBody accounts account)
	{
		accounts acc = a.findUser(account.getUsername());
		return acc != null && acc.getPassword().equals(account.getPassword());
	}

	@PostMapping("/add_profile_picture")
	public RedirectView addProfPic(@RequestParam("upload") MultipartFile file,
								   @CookieValue("username") String usr) throws IOException{

//		// No such problem
//		if(problemId.length() == 0 || !Paths.get("/mnt/problems", problemId).toFile().exists() || file.isEmpty()) {
//			return new RedirectView("/formUpload?invalid=true");
//		}
//
//		Submission submission = new Submission(problemId, file.getOriginalFilename(), username, type);
//		submission = submissions.save(submission);
//		String id = "" + submission.getSubmissionId();
//
//		File f = Paths.get("/mnt/submissions/", id).toFile();
//		if(!f.exists())
//		{
//			f.mkdirs();
//		}
//
//		File fNew = submission.getFile();
//		fNew.createNewFile();
//		FileOutputStream fout = new FileOutputStream(fNew);
//		fout.write(file.getBytes());
//		fout.close();
//
//		return new RedirectView("/results/" + submission.getSubmissionId());
		return new RedirectView("/profile");
	}
}
