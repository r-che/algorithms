package rbtree

import (
	"fmt"
	"testing"
)


func TestStringFilled(t *testing.T) {
	//nolint:lll // Predefined colored tree
	want := fmt.Sprintf(
`                                                                %[4]s%[2]s%[5]s 25                                                                                                    
                                          _____________________/     \_______________________________________________________________                                   
                                         /                                                                                           \                                  
                                    %[4]s%[2]s%[5]s 20                                                                                              %[3]s%[1]s%[5]s 35                              
                            _______/     \_______                                                        ____________________________/     \______________              
                           /                     \                                                      /                                                 \             
                      %[3]s%[1]s%[5]s 5                         %[4]s%[2]s%[5]s 22                                             %[4]s%[2]s%[5]s 30                                                    %[4]s%[2]s%[5]s 390        
              _______/     \                     /     \                                   _______/     \_______                                   _______/     \       
             /              \                   /       \                                 /                     \                                 /              \      
        %[4]s%[2]s%[5]s 3                  %[4]s%[2]s%[5]s 10          %[3]s%[1]s%[5]s 21          %[3]s%[1]s%[5]s 23                        %[3]s%[1]s%[5]s 28                        %[3]s%[1]s%[5]s 32                        %[4]s%[2]s%[5]s 37                 %[4]s%[2]s%[5]s 400 
       /     \                                                                      /     \                     /     \_______                    \                     
      /       \                                                                    /       \                   /              \                    \                    
 %[3]s%[1]s%[5]s 2           %[3]s%[1]s%[5]s 4                                                            %[4]s%[2]s%[5]s 27          %[4]s%[2]s%[5]s 29          %[4]s%[2]s%[5]s 31                 %[4]s%[2]s%[5]s 34                 %[3]s%[1]s%[5]s 38                
                                                                             /                                                /                                         
                                                                            /                                                /                                          
                                                                       %[3]s%[1]s%[5]s 26                                             %[3]s%[1]s%[5]s 33                                            
`,
	CircleRed, CircleBlack, TermRed, TermBlack, TermRst)

	// Make real tree
	tree := NewRBTree()
	for _, k := range []KeyType{
		20, 10, 30, 5, 25, 35, 37, 34, 2, 23, 27, 21, 31,
		3, 4, 28, 29, 400, 390, 38, 26, 22, 31, 32, 33,
	} {
		tree.Insert(NewRBNode(k, nil))
	}

	// Compare
	if tStr := tree.String(); tStr != want {
		t.Errorf("RBTree.String() returned:\n---\n%s\n---\nWant:\n---\n%s\n---\n", tStr, want)
	}
}

func TestStringEmpty(t *testing.T) {
	tree := NewRBTree()
	if tStr := tree.String(); tStr != strEmptyTree {
		t.Errorf("RBTree.String() returned:\n---\n%s\n---\nWant:\n---\n%s\n---\n", tStr, strEmptyTree)
	}
}
