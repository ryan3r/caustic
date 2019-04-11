package com.ryan3r.caustic;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class accountsController {
	accounts a = new accounts();
	
	@RequestMapping(method=RequestMethod.POST, value="/accounts")
	public void addAccount(@RequestBody accounts account) {
		System.out.println(account.getUsername());
	}

}
