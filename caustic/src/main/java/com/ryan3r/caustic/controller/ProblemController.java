package com.ryan3r.caustic.controller;

import com.ryan3r.caustic.repository.ProblemRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

import java.util.NoSuchElementException;

@RestController
public class ProblemController {
    @Autowired
    ProblemRepository problemRepository;

    @RequestMapping(method = RequestMethod.GET, value = "/pdf/{id}.pdf", produces = MediaType.APPLICATION_PDF_VALUE)
    @ResponseBody
    public byte[] problemPdf(@PathVariable("id") String id) {
        try {
            return problemRepository.findById(Long.parseLong(id)).get().getContent();
        } catch(NoSuchElementException e) {
            return new byte[0];
        }
    }
}