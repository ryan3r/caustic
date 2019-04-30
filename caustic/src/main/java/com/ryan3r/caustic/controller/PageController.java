package com.ryan3r.caustic.controller;

import com.ryan3r.caustic.model.Submission;
import com.ryan3r.caustic.model.accounts;
import com.ryan3r.caustic.repository.LanguageRepository;
import com.ryan3r.caustic.repository.SubmissionRepository;
import com.ryan3r.caustic.repository.accountsRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.CookieValue;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestParam;

import java.util.NoSuchElementException;

@Controller
public class PageController {

    @Value("${spring.application.name}")
    String appName;

    private final accountsRepository accRep;

    @Autowired
    private SubmissionRepository submissions;
    @Autowired
    private LanguageRepository languageRepository;

    @Autowired
    public PageController(accountsRepository accRep){
        this.accRep = accRep;
    }

    @GetMapping(path="/")
    public String index(Model model){
        model.addAttribute("appName", appName);
        return "index";
    }

    @GetMapping("/profile")
    public String profile(@CookieValue("username") String username, Model model) {

        accounts account = accRep.findUser(username);
        model.addAttribute("account", account);
        return "profile";
    }

    @GetMapping("/results/{id}")
    public String results(@PathVariable("id") String id, Model model) {
        Submission submit = null;
        try {
            submit = submissions.findById(Long.parseLong(id)).get();
        } catch(NoSuchElementException ex) {
            model.addAttribute("error", "Submission not found");
        }
        model.addAttribute("result", submit);
        return "results";
    }

    @GetMapping("/scoreboard")
    public String scoreboard(Model model){
        model.addAttribute("scores", submissions.findAll());
        return "scoreboard";
    }

    @GetMapping("/accountSetup")
    public String accountSetup(Model model){
        return "accountSetup";
    }

    @GetMapping("/login-newaccount")
    public String login(Model model){
        return "login";
    }

    @GetMapping("/formUpload")
    public String formUpload(@RequestParam(value = "invalid", required = false) boolean invalid, Model model){
        model.addAttribute("invalid", invalid);
        model.addAttribute("types", languageRepository.findAll());
        return "formUpload";
    }
}