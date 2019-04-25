package com.ryan3r.caustic.model;

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
     * @param className The name to use when running a file
     */
    public Submission(String _problem, String _fileName) {
        problem = _problem;
        fileName = _fileName;
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
}
