package com.ryan3r.caustic.controller;

import com.ryan3r.caustic.model.accounts;
import com.ryan3r.caustic.repository.accountsRepository;
import com.ryan3r.caustic.repository.submitsRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

@Controller
public class PageController {

    @Value("${spring.application.name}")
    String appName;

    private final accountsRepository accRep;
    private final submitsRepository subRep;

    @Autowired
    public PageController(accountsRepository accRep, submitsRepository subRep){
        this.accRep = accRep;
        this.subRep = subRep;
    }

    @GetMapping(path="/")
    public String index(Model model){
        model.addAttribute("appName", appName);
        return "index";
    }

    //TODO: Change this to cookie, and save user as cookie if we want user to see their page
    @GetMapping("/profile/{username}")
    public String profile(@PathVariable("username") String username, Model model) {

        accounts account = accRep.findUser(username);
        model.addAttribute("account", account);
        return "profile";
    }

    //TODO: Pull from Submission or submits
    @GetMapping("/results")
    public String results(Model model) {
        //model.addAttribute("result", result);
        return "results";
    }

    //TODO: Is this right? Submission or submits?
    @GetMapping("/scoreboard")
    public String scoreboard(Model model){
        model.addAttribute("scores", subRep.findAll());
        return "scoreboard";
    }

    @GetMapping("/accountSetup")
    public String accountSetup(Model model){
        return "accountSetup";
    }

    @GetMapping("/login-newaccount")
    public String login(Model model){
        return "home";
    }

    @GetMapping("/formUpload")
    public String formUpload(Model model){
        return "formUpload";
    }
}