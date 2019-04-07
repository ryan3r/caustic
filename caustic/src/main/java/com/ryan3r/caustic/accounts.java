package com.ryan3r.caustic;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
public class accounts {
	@Id
	@GeneratedValue(strategy=GenerationType.AUTO) private Long id;
	private String username;
	private String password;
	
	protected accounts() {}
	
	public accounts(String username, String password)
	{
		this.username = username;
		this.password = password;
	}
	
	@Override
	public String toString() 
	{
		return String.format("account[id=%d, username='%s', password='%s']", 
				id, username, password);
	}
}
