package com.ryan3r.caustic.model;

import java.io.File;
import java.nio.file.Files;
import java.nio.file.InvalidPathException;
import java.nio.file.Paths;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

import org.springframework.lang.NonNull;

@Entity
public class Submission {
    public Submission() {}

    /**
     * Create a new submission
     * @param problem The name of the problem this submission solves
     * @param fileName The name of the submission file in /mnt/submissions
     */
    public Submission(String _problem, String _fileName) {
        problem = _problem;
        fileName = _fileName;
        status = SubmissionStatus.NEW;
    }

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    long submissionId;

    @NonNull
    SubmissionStatus status;

    @NonNull
    String fileName;

    @NonNull
    String problem;

    // The status/result of a submission
    public enum SubmissionStatus {
        NEW,
        RUNNING,
        COMPILE_ERROR,
        OK,
        WRONG,
        TIME_LIMIT,
        EXCEPTION
    }

    /**
     * @return the submissionId
     */
    public long getSubmissionId() {
        return submissionId;
    }

    /**
     * @param submissionId the submissionId to set
     */
    public void setSubmissionId(long submissionId) {
        this.submissionId = submissionId;
    }

    /**
     * @return the status
     */
    public SubmissionStatus getStatus() {
        return status;
    }

    /**
     * @param status the status to set
     */
    public void setStatus(SubmissionStatus status) {
        this.status = status;
    }

    /**
     * @return the fileName
     */
    public String getFileName() {
        return fileName;
    }

    /**
     * @param fileName the fileName to set
     */
    public void setFileName(String fileName) {
        this.fileName = fileName;
    }

    /**
     * @return the problem
     */
    public String getProblem() {
        return problem;
    }

    /**
     * @param problem the problem to set
     */
    public void setProblem(String problem) {
        this.problem = problem;
    }

    /**
     * Get the submission the user uploaded
     */
    public File getFile() {
        try {
            return Paths.get("/mnt/submissions/", submissionId + "", fileName).toFile();
        } catch(InvalidPathException ex) {
            ex.printStackTrace();
            return null;
        }
    }

    /**
     * Get the content of the submission
     * @return
     */
    public String getContent() {
        try {
            return new String(Files.readAllBytes(Paths.get("/mnt/submissions/", submissionId + "", fileName)));
        } catch(Exception ex) {
            ex.printStackTrace();
            return "";
        }
    }

    /**
     * Submission status as a string
     */
    public String getStatusText() {
        switch(status) {
            case NEW:
                return "New";
            case TIME_LIMIT:
                return "Time limit exceded";
            case COMPILE_ERROR:
                return "Compile Error";
            case EXCEPTION:
                return "Runtime Error";
            case RUNNING:
                return "Running";
            case OK:
                return "Accepted";
            default:
                return "Wrong Answer";
        }
    }

    /**
     * Get a wrong, right or loading image
     * @return
     */
    public String getStatusImage() {
        switch(status) {
            case NEW:
            case RUNNING:
                return "/loading.gif";
            case OK:
                return "/correct.png";
            default:
                return "/wrong.png";
        }
    }
}
