// TreeParityTest.java - 验证Go实现与Java 1.12树木生成的数学一致性

import java.util.Random;

/**
 * TreeParityTest - 验证Go实现与Java 1.12树木生成的数学一致性
 * 
 * 测试项目:
 * 1. 深色橡木高度公式: nextInt(3) + nextInt(2) + 6
 * 2. 橡木角落逻辑: Math.abs(l1) != j1 || Math.abs(j2) != j1 || (rand.nextInt(2) != 0
 * && i4 != 0)
 * 3. 半径计算: 1 - i4/2 (整数除法)
 */
public class TreeParityTest {

    public static void main(String[] args) {
        long seed = 12345L;
        Random rand = new Random(seed);

        System.out.println("=== Java 1.12 Tree Generation Parity Test ===");
        System.out.println("Seed: " + seed);
        System.out.println();

        // Test 1: Dark Oak Height Formula
        System.out.println("--- Test 1: Dark Oak Height Formula ---");
        System.out.println("Formula: nextInt(3) + nextInt(2) + 6");
        rand = new Random(seed);
        for (int i = 0; i < 20; i++) {
            int height = rand.nextInt(3) + rand.nextInt(2) + 6;
            System.out.println("Tree " + i + ": height=" + height);
        }
        System.out.println();

        // Test 2: Oak Leaf Radius Calculation (Integer Division)
        System.out.println("--- Test 2: Leaf Radius (Integer Division) ---");
        System.out.println("Formula: j1 = 1 - i4/2");
        for (int i4 = -3; i4 <= 0; i4++) {
            int j1 = 1 - i4 / 2;
            System.out.println("i4=" + i4 + " -> j1=" + j1);
        }
        System.out.println();

        // Test 3: Corner Rounding Logic
        System.out.println("--- Test 3: Corner Rounding Logic ---");
        System.out.println("Formula: place if (abs(l1)!=j1 || abs(j2)!=j1 || (nextInt(2)!=0 && i4!=0))");
        rand = new Random(seed);
        int treeHeight = 5;
        int placedCount = 0;
        int skippedCount = 0;

        for (int i3 = 0 - 3 + treeHeight; i3 <= 0 + treeHeight; i3++) {
            int i4 = i3 - (0 + treeHeight);
            int j1 = 1 - i4 / 2;

            for (int k1 = 0 - j1; k1 <= 0 + j1; k1++) {
                int l1 = k1 - 0;

                for (int i2 = 0 - j1; i2 <= 0 + j1; i2++) {
                    int j2 = i2 - 0;

                    // The actual Java 1.12 condition
                    boolean shouldPlace = Math.abs(l1) != j1 || Math.abs(j2) != j1 ||
                            (rand.nextInt(2) != 0 && i4 != 0);

                    if (shouldPlace) {
                        placedCount++;
                    } else {
                        skippedCount++;
                        System.out.println("SKIP: i4=" + i4 + " j1=" + j1 + " l1=" + l1 + " j2=" + j2);
                    }
                }
            }
        }
        System.out.println("Total placed: " + placedCount + ", skipped: " + skippedCount);
        System.out.println();

        // Test 4: RNG Sequence for verification
        System.out.println("--- Test 4: RNG Sequence (first 50 values) ---");
        rand = new Random(seed);
        System.out.print("nextInt(100): ");
        for (int i = 0; i < 10; i++) {
            System.out.print(rand.nextInt(100) + " ");
        }
        System.out.println();

        rand = new Random(seed);
        System.out.print("nextBool: ");
        for (int i = 0; i < 10; i++) {
            System.out.print(rand.nextBoolean() + " ");
        }
        System.out.println();

        System.out.println();
        System.out.println("=== Test Complete ===");
        System.out.println("Run the corresponding Go test and compare outputs.");
    }
}
