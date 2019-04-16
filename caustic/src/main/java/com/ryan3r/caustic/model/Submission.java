package com.ryan3r.caustic.model;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

import org.springframework.lang.NonNull;

@Entity
public class Submission {
    Submission() {}

    /**
     * Create a new submission
     * @param problem The id of the problem this submission solves
     * @param fileName The name of the submission file in /mnt/submissions
     * @param className The name to use when running a file
     */
    Submission(long _problem, String _fileName, String _className) {
        problem = _problem;
        fileName = _fileName;
        className = _className;
    }

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    long submissionId;

    @NonNull
    SubmissionStatus status;

    @NonNull
    String fileName;

    @NonNull
    long problem;

    // The percent of inputs that have been tested (out of 100)
    @NonNull
    int progress; 

    String className;

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
}