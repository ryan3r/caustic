package com.ryan3r.caustic;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class accountsController {
	@Autowired
	accountsRepository a;
	
	@RequestMapping(method=RequestMethod.POST, value="/accounts")
	public String addAccount(@RequestBody accounts account) {
		a.save(account);
		return "Yeet";
	}
}
