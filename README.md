gof
===

GOF: The golang mvc web framework 

View:
  1. You can declare the html file(.gohtml) mix with html&golang. eg:view/home/index.gohtml
    
  2. Then,use the goftool,you can compile the html files(.gohtml) into golang files.  
  such as the html file path is "view/home/index.gohtml",
  compile the view directory in command:                                                     
  
  ./gof compileview view/ view/  view/helper.go

  the golang file will be generated:view/V_home_index.go . If any error,the compiler will alert

  3.The html result will render at runtime in high performance(since the code is purely golang) 


Controller:

Model:
