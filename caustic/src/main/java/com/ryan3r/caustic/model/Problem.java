package com.ryan3r.caustic.model;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import org.springframework.lang.NonNull;

@Entity
public class Problem {
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    long id;

    @NonNull
    String name;

    @NonNull
    String pdfPath;

    @Override
    public String toString() {
        return name;
    }

    /**
     * @return the id
     */
    public long getId() {
        return id;
    }

    /**
     * @param id the id to set
     */
    public void setId(long id) {
        this.id = id;
    }

    /**
     * @return the name
     */
    public String getName() {
        return name;
    }

    /**
     * @param name the name to set
     */
    public void setName(String name) {
        this.name = name;
    }

    /**
     * @return the pdfPath
     */
    public String getPdfPath() {
        return pdfPath;
    }

    /**
     * @param pdfPath the pdfPath to set
     */
    public void setPdfPath(String pdfPath) {
        this.pdfPath = pdfPath;
    }

    /**
     * Get the content of the submission file
     */
    public byte[] getContent() {
        try {
            return Files.readAllBytes(Paths.get(pdfPath));
        } catch(IOException e) {
            return new byte[0];
        }
    }

    /**
     * Get the url for the pdf
     */
    public String getPdfUrl() {
        return "/pdf/" + id + ".pdf";
    }

    /**
     * Get the url for the submission page
     */
    public String getUploadUrl() {
        return "/formUpload/" + id;
    }
}
