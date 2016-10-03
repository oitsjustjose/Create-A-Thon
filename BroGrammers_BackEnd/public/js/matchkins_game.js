(function(){

	'use strict';

	var Memory = {

		init: function(cards){
			this.$game = $(".game");
			this.$modal = $(".modal");
			this.$overlay = $(".modal-overlay");
			this.$restartButton = $("button.restart");
			this.cardsArray = $.merge(cards, cards);
			this.shuffleCards(this.cardsArray);
			this.setup();
		},

		shuffleCards: function(cardsArray){
			this.$cards = $(this.shuffle(this.cardsArray));
		},

		setup: function(){
			this.html = this.buildHTML();
			this.$game.html(this.html);
			this.$memoryCards = $(".card");
			this.binding();
			this.paused = false;
     	this.guess = null;
		},

		binding: function(){
			this.$memoryCards.on("click", this.cardClicked);
			this.$restartButton.on("click", $.proxy(this.reset, this));
		},
		// kinda messy but hey
		cardClicked: function(){
			var _ = Memory;
			var $card = $(this);
			if(!_.paused && !$card.find(".inside").hasClass("matched") && !$card.find(".inside").hasClass("picked")){
				$card.find(".inside").addClass("picked");
				if(!_.guess){
					_.guess = $(this).attr("data-id");
				} else if(_.guess == $(this).attr("data-id") && !$(this).hasClass("picked")){
					$(".picked").addClass("matched");
					_.guess = null;
				} else {
					_.guess = null;
					_.paused = true;
					setTimeout(function(){
						$(".picked").removeClass("picked");
						Memory.paused = false;
					}, 600);
				}
				if($(".matched").length == $(".card").length){
					_.win();
				}
			}
		},

		win: function(){
			this.paused = true;
			setTimeout(function(){
				Memory.showModal();
				Memory.$game.fadeOut();
			}, 1000);
		},

		showModal: function(){
			this.$overlay.show();
			this.$modal.fadeIn("slow");
		},

		hideModal: function(){
			this.$overlay.hide();
			this.$modal.hide();
		},

		reset: function(){
			this.hideModal();
			this.shuffleCards(this.cardsArray);
			this.setup();
			this.$game.show("slow");
		},

		// Fisher--Yates Algorithm -- http://bost.ocks.org/mike/shuffle/
		shuffle: function(array){
			var counter = array.length, temp, index;
	   	// While there are elements in the array
	   	while (counter > 0) {
        	// Pick a random index
        	index = Math.floor(Math.random() * counter);
        	// Decrease counter by 1
        	counter--;
        	// And swap the last element with it
        	temp = array[counter];
        	array[counter] = array[index];
        	array[index] = temp;
	    	}
	    	return array;
		},

		buildHTML: function(){
			var frag = '';
			this.$cards.each(function(k, v){
				frag += '<div class="card" data-id="'+ v.id +'"><div class="inside">\
				<div class="front"><img src="'+ v.img +'"\
				alt="'+ v.name +'" /></div>\
				<div class="back"><img src="http://puu.sh/ruvNj/88d32a53a9.png"\
				alt="Codepen" /></div></div>\
				</div>';
			});
			return frag;
		}
	};

	var cards = [
		{
			name: "apple_blossom",
			img: "http://api.shopkinsworld.com/media/SPKS1/SPK_001.png",
			id: 1,
		},
		{
			name: "rockin_broc",
			img: "http://api.shopkinsworld.com/media/SPKS1/SPK_002.png",
			id: 2
		},
		{
			name: "strawberry_kiss",
			img: "http://api.shopkinsworld.com/media/SPKS1/SPK_003.png",
			id: 3
		},
		{
			name: "pineapple_crush",
			img: "http://api.shopkinsworld.com/media/SPKS1/SPK_004.png",
			id: 4
		},
		{
			name: "melonie_pips",
			img: "http://api.shopkinsworld.com/media/SPKS1/SPK_005.png",
			id: 5
		},
		{
			name: "miss_mushy-moo",
			img: "http://api.shopkinsworld.com/media/SPKS1/SPK_006.png",
			id: 6
		},
		{
			name: "posh_pear",
			img: "http://api.shopkinsworld.com/media/SPKS1/SPK_007.png",
			id: 7
		},
		{
			name: "apple_blossom_green",
			img: "http://api.shopkinsworld.com/media/SPKS1/SPK_008.png",
			id: 8
		},
		{
			name: "rockin_broc_blue",
			img: "http://api.shopkinsworld.com/media/SPKS1/SPK_009.png",
			id: 9
		},
		{
			name: "strawberry_kiss_green",
			img: "http://api.shopkinsworld.com/media/SPKS1/SPK_010.png",
			id: 10
		},
		{
			name: "pineapple_crush_pink",
			img: "http://api.shopkinsworld.com/media/SPKS1/SPK_011.png",
			id: 11
		},
		{
			name: "melonie_pips_honeydew",
			img: "http://api.shopkinsworld.com/media/SPKS1/SPK_012.png",
			id: 12
		},
	];

	Memory.init(cards);


})();
