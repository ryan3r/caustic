package com.ryan3r.caustic.controller;

import com.ryan3r.caustic.repository.accountsRepository;
import com.ryan3r.caustic.model.Scoreboard;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.CookieValue;
import org.springframework.web.bind.annotation.GetMapping;

import java.util.ArrayList;
import java.util.List;

@Controller
public class PageController {

    @Value("${spring.application.name}")
    String appName;

    @GetMapping(path="/")
    public String index(Model model){
        model.addAttribute("appName", appName);
        return "index";
    }

    accountsRepository a;

    @GetMapping("/profile")
    public String profile(Model model/*,
                          @CookieValue("username") String usrNameCookie,
                          @CookieValue("accType") String accType*/
                         ) {

        //model.addAttribute("username", usrNameCookie);
        //model.addAttribute("accType", accType);
        return "profile";
    }

    @GetMapping("/results")
    public String results(Model model) {
        //model.addAttribute("result", result);
        return "results";
    }

    private List<Scoreboard> scoreList = new ArrayList<>();

    @GetMapping("/scoreboard")
    public String scoreboard(Model model){
        model.addAttribute("scores", scoreList);
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