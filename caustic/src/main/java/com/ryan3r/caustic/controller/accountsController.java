package com.ryan3r.caustic.controller;

import com.ryan3r.caustic.model.accounts;
import com.ryan3r.caustic.repository.accountsRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.servlet.view.RedirectView;

import java.io.IOException;

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


	//TODO: create this folder in the DB so that we can save profile pics
//	@PostMapping("/add_profile_picture")
//	public RedirectView addProfPic(@RequestParam("upload") MultipartFile file,
//								   @CookieValue("username") String usrname) throws IOException {
//
//		// Didn't upload a file, or uploaded an empty file.
//		// Note: To reach this page a user must be logged in, so usr==null is never true, but we'll check just in case.
////		if(usrname==null || file.isEmpty()) {
////			return new RedirectView("/profile?error=2");
////		}
////
////		accounts acc = a.findUser(usrname);
////		a.insertProfPic(usrname);
////
////		Path filepath = Paths.get("/mnt/profPics", usrname);
////		if(filepath == null)
////		{
////			return new RedirectView("/profile?error=2");
////		}
////
////		acc.setPathToProfPic("/mnt/profPics/"+usrname);
////		file.transferTo(filepath);
//
//		return new RedirectView("/profile");
//	}


}
