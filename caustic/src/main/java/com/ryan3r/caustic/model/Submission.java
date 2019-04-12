package com.ryan3r.caustic.model;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

import org.springframework.lang.NonNull;

@Entity
public class Submission {
    Submission() {

    }

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    long submissionId;

    @NonNull
    SubmissionStatus status;

    @NonNull
    String fileName;

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