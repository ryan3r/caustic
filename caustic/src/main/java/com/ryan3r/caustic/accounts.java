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
	private String accType;
	
	public accounts() {}
	
	public accounts(String username, String password, String accType)
	{
		this.username = username;
		this.password = password;
		this.accType = accType;
	}
	
	@Override
	public String toString() 
	{
		return String.format("account[id=%d, username='%s', password='%s']", 
				id, username, password);
	}
	
	public void setUsername(String username)
	{
		this.username = username;
	}
	
	public void setPassword(String password)
	{
		this.password = password;
	}
	
	public void setAccType(String accType)
	{
		this.accType = accType;
	}
	
	public String getUsername()
	{
		return this.username;
	}
	
	public String getPassword()
	{
		return this.password;
	}
	
	public String getAccType()
	{
		return this.accType;
	}
}
