package main.java.usermgmt;

import io.codenotary.immudb4j.FileImmuStateHolder;
import io.codenotary.immudb4j.ImmuClient;
import io.codenotary.immudb4j.user.Permission;
import io.codenotary.immudb4j.user.User;

import java.io.IOException;
import java.util.List;


public class UserMgmtSample {

    static final String USERNAME = "immudb";
    static final String PASSWORD = "immudb";
    static final String DATABASE = "defaultdb";

    public static void main(String[] args) {

        ImmuClient immuClient = null;


        FileImmuStateHolder stateHolder = null;
        try {
            stateHolder = FileImmuStateHolder.newBuilder()
                    .setStatesFolder("states")
                    .build();
        } catch (IOException e) {
            System.err.println(">>> Failed to create init stateHolder. Reason: " + e.getMessage());
        }

        immuClient = ImmuClient.newBuilder()
                .setStateHolder(stateHolder)
                .setServerUrl("localhost")
                .setServerPort(3322)
                .setWithAuthToken(true)
                .build();

        immuClient.login(USERNAME, PASSWORD);

        immuClient.useDatabase(DATABASE);

        // ---------- doing/testing some logic ----------

        String username = "dxps";
        String password = "paSs123$%^";
        immuClient.createUser(username, password, Permission.PERMISSION_RW, DATABASE);

        System.out.println(">>> listUsers:");
        List<User> users = immuClient.listUsers();
        users.forEach(user -> System.out.println("\t- " + user));

        // ----------------------------------------------

        immuClient.logout();
        immuClient.shutdown();
    }

}
