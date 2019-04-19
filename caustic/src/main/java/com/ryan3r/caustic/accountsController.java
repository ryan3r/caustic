package com.ryan3r.caustic;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

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
		if(a.findById(account.getUsername()) == null)
			return false;
		else if(a.findById(account.getUsername()).get().getPassword() != account.getPassword())
			return false;
		return true;
	}
}
