package dev.haruki7049.Spelling;

import com.almasb.fxgl.app.GameApplication;
import com.almasb.fxgl.app.GameSettings;

public class Main extends GameApplication {
    @Override
    protected void initSettings(GameSettings settings) {
        settings.setWidth(800);
        settings.setHeight(600);
        settings.setTitle("Spelling");
    }

    @Override
    protected void initGame() {
    }

    public static void main(String[] args) {
        launch(args);
    }
}
