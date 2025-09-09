package main

import "testing"
import "fmt"

type testCase struct {
    a, b     int
    expected int
}
type testCaseFloat struct {
	a, b, expected float64
}

func TestFirstWay(t *testing.T) {

    testCases := []testCase{
        {5, 8, 40},     
        {0, 8, 0},      
        {5, 0, 0},      
        {-5, 8, -40},  
        {5, -8, -40},   
        {-5, -8, 40},  
    }

    for i, tc := range testCases {
        t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
            actual, err := firstWay(tc.a, tc.b)
            if err != nil {
                t.Errorf("Should not produce an error: %v", err)
            }
            if tc.expected != actual {
                t.Errorf("For a=%d, b=%d: expected %d, got %d", tc.a, tc.b, tc.expected, actual)
            }
        })
    }
}
func TestSecondWay(t *testing.T) {

   testCase := []testCaseFloat{
        {5, 8, 40},     
        {0, 8, 0},      
        {5, 0, 0},   
        {-5, 8, -40},  
        {5, -8, -40},   
        {-5, -8, 40},  
    }

	for i, tc := range testCase {
        t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
            actual, err := secondWay(tc.a, tc.b)
            if err != nil {
                t.Errorf("Should not produce an error: %v", err)
            }
            if tc.expected != actual {
                t.Errorf("For a=%f, b=%f: expected %f, got %f", tc.a, tc.b, tc.expected, actual)
            }
        })
    }
}
func TestThirdWay(t *testing.T){

	testCase := []testCaseFloat{
        {5, 8, 40},     
        {0, 8, 0},      
        {5, 0, 0},   
        {-5, 8, -40},  
        {5, -8, -40},   
        {-5, -8, 40},  
	}

	for i, tc := range testCase {
        t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
            actual, err := thirdWay(tc.a, tc.b)
            if err != nil {
                t.Errorf("Should not produce an error: %v", err)
            }
            if tc.expected != actual {
                t.Errorf("For a=%f, b=%f: expected %f, got %f", tc.a, tc.b, tc.expected, actual)
            }
        })
    }
}
func TestFourthWay(t *testing.T){

	testCase := []testCaseFloat{
        {5, 8, 40},     
        {0, 8, 0},      
        {5, 0, 0},   
        {-5, 8, -40},  
        {5, -8, -40},   
        {-5, -8, 40},  
	}

	for i, tc := range testCase {
        t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
            actual, err := fourthWay(tc.a, tc.b)
            if err != nil {
                t.Errorf("Should not produce an error: %v", err)
            }
            if tc.expected != actual {
                t.Errorf("For a=%f, b=%f: expected %f, got %f", tc.a, tc.b, tc.expected, actual)
            }
        })
    }
}
func TestFifthWay(t *testing.T) {

   testCase := []testCase{
        {5, 8, 40},     
        {0, 8, 0},      
        {5, 0, 0},   
        {-5, 8, -40},  
        {5, -8, -40},   
        {-5, -8, 40},  
    }

	for i, tc := range testCase {
        t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
            actual, err := fifthWay(tc.a, tc.b)
            if err != nil {
                t.Errorf("Should not produce an error: %v", err)
            }
            if tc.expected != actual {
                t.Errorf("For a=%d, b=%d: expected %d, got %d", tc.a, tc.b, tc.expected, actual)
            }
        })
    }
}

func TestSixthWay(t *testing.T) {

   testCase := []testCase{
        {5, 8, 40},     
        {0, 8, 0},      
        {5, 0, 0},   
        {-5, 8, -40},  
        {5, -8, -40},   
        {-5, -8, 40},
        {1, 1, 1}, 
        {1, -1, -1}, 
        {0, 0, 0},  
        {5, 1, 5},
        {1, 5, 5},
    }

	for i, tc := range testCase {
        t.Run(fmt.Sprintf("Case%d", i+1), func(t *testing.T) {
            actual, err := sixthWay(tc.a, tc.b)
            if err != nil {
                t.Errorf("Should not produce an error: %v", err)
            }
            if tc.expected != actual {
                t.Errorf("For a=%d, b=%d: expected %d, got %d", tc.a, tc.b, tc.expected, actual)
            }
        })
    }
}