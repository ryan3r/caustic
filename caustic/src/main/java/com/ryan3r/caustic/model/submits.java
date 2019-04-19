package com.ryan3r.caustic.model;
import java.io.File;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
public class submits {
	@Id
    @GeneratedValue(strategy=GenerationType.AUTO) private Long id;
	private File file;
	private String submissionID;
	
	public submits() {}
	
	public submits(String submitID, File f)
	{
		this.submissionID = submitID;
		this.file = f;
	}
	
	public File getFile()
	{
		return this.file;
	}
	
	public String getSubmissionID()
	{
		return this.submissionID;
	}
	
	public void setFile(File f)
	{
		this.file = f;
	}
	
	public void setSumbissionID(String submitId)
	{
		this.submissionID = submitId;
	}
	
}
