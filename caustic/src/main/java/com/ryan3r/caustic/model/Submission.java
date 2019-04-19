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
     * @param problem The name of the problem this submission solves
     * @param fileName The name of the submission file in /mnt/submissions
     * @param className The name to use when running a file
     */
    Submission(String _problem, String _fileName) {
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
}
