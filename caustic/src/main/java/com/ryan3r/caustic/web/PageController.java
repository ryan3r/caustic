package com.ryan3r.caustic.web;

import com.ryan3r.caustic.model.Scoreboard;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;

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

    @GetMapping("/profile")
    public String profile(Model model,
                          @RequestParam(value="firstName", required=false, defaultValue="Zach") String firstName,
                          @RequestParam(value="lastName", required=false, defaultValue="Gorman") String lastName,
                          @RequestParam(value="email", required=false, defaultValue="zgorman2@iastate.edu") String email
                         ) {

        model.addAttribute("firstName", firstName);
        model.addAttribute("lastName", lastName);
        model.addAttribute("email", email);

        return "profile";
    }

    @GetMapping("/results")
    public String results(Model model/*, @RequestParam(value="result", required=false, defaultValue="fail") String result*/) {
        /*model.addAttribute("result", result);*/
        return "results";
    }

    List<Scoreboard> scoreList = new ArrayList<>();

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