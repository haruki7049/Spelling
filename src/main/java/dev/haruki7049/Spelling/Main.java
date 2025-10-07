package dev.haruki7049.Spelling;

import static com.almasb.fxgl.dsl.FXGL.*;

import com.almasb.fxgl.app.GameApplication;
import com.almasb.fxgl.app.GameSettings;
import com.almasb.fxgl.input.KeyTrigger;
import com.almasb.fxgl.input.TriggerListener;
import com.almasb.fxgl.logging.Logger;
import javafx.beans.property.StringProperty;

public class Main extends GameApplication {

  private StringProperty currentInput;
  private static final String ALLOWED_CHARS =
      "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 -',.!?";
  private static final int MAX_LENGTH = 20;
  private final Logger log = Logger.get(Main.class);

  @Override
  protected void initSettings(GameSettings settings) {
    settings.setWidth(800);
    settings.setHeight(600);
    settings.setTitle("Spelling");
  }

  @Override
  protected void initInput() {
    getInput()
        .addTriggerListener(
            new TriggerListener() {
              @Override
              protected void onKeyBegin(KeyTrigger keyTrigger) {
                var code = keyTrigger.getKey();

                boolean shift = getInput().isHeld(javafx.scene.input.KeyCode.SHIFT);

                // 英数字キーのみ判定
                if (code.isLetterKey()) {
                  char c = code.getName().charAt(0);
                  if (!shift) c = Character.toLowerCase(c);
                  if (ALLOWED_CHARS.indexOf(c) != -1) {
                    log.info("Accepted: " + c);
                  }
                } else if (code.isDigitKey()) {
                  char c = code.getName().charAt(0);
                  if (ALLOWED_CHARS.indexOf(c) != -1) {
                    log.info("Accepted: " + c);
                  }
                }

                // 記号・スペース類
                else {
                  switch (code) {
                    case SPACE -> logIfAllowed(" ");
                    case COMMA -> logIfAllowed(",");
                    case PERIOD -> logIfAllowed(".");
                    case QUOTE -> logIfAllowed("'");
                    case MINUS -> logIfAllowed("-");
                    case DIGIT1 -> {
                      if (shift) logIfAllowed("!");
                    }
                    case SLASH -> {
                      if (shift) logIfAllowed("?");
                    }
                  }
                }
              }

              private void logIfAllowed(String s) {
                if (ALLOWED_CHARS.contains(s)) log.info("Accepted: " + s);
              }
            });
  }

  public static void main(String[] args) {
    launch(args);
  }
}
