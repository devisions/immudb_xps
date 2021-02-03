package main.java.io.codenotary.immudb.helloworld;

import io.codenotary.immudb4j.FileRootHolder;
import io.codenotary.immudb4j.ImmuClient;
import io.codenotary.immudb4j.crypto.VerificationException;

import java.io.IOException;

public class App {

    static final String USERNAME = "immudb";
    static final String PASSWORD = "immudb";
    static final String DATABASE = "defaultdb";

    public static void main(String[] args) {

        ImmuClient client = null;

        try {

            FileRootHolder rootHolder = FileRootHolder.newBuilder()
                    .setRootsFolder("./helloworld_immudb_roots").build();

            client = ImmuClient.newBuilder()
                    .setServerUrl("localhost")
                    .setServerPort(3322)
                    .setRootHolder(rootHolder).build();

            client.login(USERNAME, PASSWORD);

            client.useDatabase(DATABASE);

            client.set("hello", "immutable world!".getBytes());

            byte[] v = client.safeGet("hello");

            System.out.format("(%s, %s)%n", "hello", new String(v));


        } catch (IOException e) {
            e.printStackTrace();
        } catch (VerificationException e) {
            // Tampering detected!
            // This means the history of changes has been tampered.
            e.printStackTrace();
            System.exit(1);
        } finally {
            if (client != null) {
                client.logout();
            }
        }

    }

}
